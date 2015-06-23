package parsing

import "text/scanner"
import "io"

var stdTMap = map[rune]string{
	scanner.Ident: "Ident",
	scanner.Int: "Int",
	scanner.Float: "Float",
	scanner.Char: "Char",
	scanner.String: "String",
	scanner.RawString: "RawString",
	scanner.Comment: "Comment",
}

type StdScanner struct{
	s    scanner.Scanner
	tmap map[rune]string
	kmap map[string]string
}
func (s *StdScanner) Init(kmap map[string]string,src io.Reader) *StdScanner {
	s.tmap = stdTMap
	s.kmap = kmap
	if s.kmap==nil { s.kmap = map[string]string{} }
	s.s.Init(src)
	return s
}

func (s *StdScanner) FirstToken() *Token {
	var f Generator
	f = Generator(func() *Token {
		r := s.s.Scan()
		if r == scanner.EOF { return Retempty() }
		t := s.s.TokenText()
		d := ""
		y := ""
		if x,ok := s.kmap[t]; ok {
			d = x
		}else if x,ok := s.tmap[r]; ok {
			y = x
		}else{
			y = t
		}
		return MakeToken(y,d,t,s.s.Pos(),f)
	})
	return f()
}


