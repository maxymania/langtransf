package syntaxfile

import "text/scanner"
import u "unicode"
import "fmt"
import "strconv"
import "io"

func firstLetter(s string) rune {
	for _,r := range s { return r }
	return 0
}

type sfparser struct{
	s scanner.Scanner
	t rune
}
func (s* sfparser) next() {
	if s.t == scanner.EOF { return }
	s.t = s.s.Scan()
}
func (s* sfparser) assert(ex string) {
	if ex==s.s.TokenText() { return }
	fmt.Println(s.s.Pos(),": expected ",ex," but got ",s.s.TokenText())
	panic("SoftPanic: syntax")
}
func (s* sfparser) unquote() string {
	t,_ := strconv.Unquote(s.s.TokenText())
	return t
}

func (s* sfparser) rule0(m *Modifier) {
	for {
		switch s.t {
		case scanner.Ident:
			ident := s.s.TokenText()
			s.next()
			if u.IsUpper(firstLetter(ident)) {
				m.Data = append(m.Data,&MatchToken{ident})
			}else{
				m.Data = append(m.Data,&CallRule{ident})
			}
		case scanner.String:
			m.Data = append(m.Data,&MatchText{s.unquote()})
			s.next()
		case '+','*','?':
			modp := len(m.Data)-1
			if modp>=0 {
				m.Data[modp] = &Modifier{[]Rule{m.Data[modp]},byte(s.t),0}
			}
			s.next()
		case '(':
			s.next()
			m.Data = append(m.Data,s.rule())
			s.assert(")")
			s.next()
		case '~':
			s.next()
			s.assert("(")
			s.next() // (
			gn := ""
			switch s.t {
			case scanner.RawString,scanner.String:
				gn = s.unquote()
				s.next()
			case ':',',':s.next()
			}
			m.Data = append(m.Data,&Group{s.rule(),gn})
			s.assert(")")
			s.next() // )
		case '#':
			s.next()
			aty := s.s.TokenText()
			atyp := len(m.Data)-1
			s.next()
			if atyp>=0 {
				switch aty{
				case "-":
					m.Data[atyp] = &Modifier{[]Rule{m.Data[atyp]},M_DROP,0}
				case "e","verbose":
					if _,ok := m.Data[atyp].(*Modifier); ok {
						m.Data[atyp].(*Modifier).Flags |= F_VERBOSE
					}
				case "E","mute":
					if _,ok := m.Data[atyp].(*Modifier); !ok {
						m.Data[atyp] = &Modifier{[]Rule{m.Data[atyp]},M_SEQ,0}
					}
					m.Data[atyp].(*Modifier).Flags |= F_MUTE
				case "p","precise":
					m.Data[atyp] = &Modifier{[]Rule{m.Data[atyp]},M_SEQ,0}
				}
			}
		default: return
		}
	}
}

func (s* sfparser) rule1(m *Modifier) {
	md := &Modifier{Mode: M_SEQ}
	s.rule0(md)
	for {
		switch s.t {
		case '|':
			s.next()
			m.Data = append(m.Data,md)
			md = &Modifier{Mode: M_SEQ}
			s.rule0(md)
		default:
			m.Data = append(m.Data,md)
			return
		}
	}
}

func (s* sfparser) rule() Rule {
	md := &Modifier{Mode: M_OR}
	s.rule1(md)
	return md
}

func (s* sfparser) parse() SyntaxFile {
	sf := make(SyntaxFile)
	for s.t!=scanner.EOF {
		t := s.s.TokenText()
		s.next()
		s.assert("=")
		s.next()
		r := s.rule()
		s.assert(";")
		s.next()
		sf[t]=r
	}
	return sf
}

// Parses an SyntaxFile. On Syntax error it panics.
func Parse(src io.Reader) SyntaxFile {
	sf := &sfparser{}
	sf.s.Init(src)
	sf.next()
	return sf.parse()
}

