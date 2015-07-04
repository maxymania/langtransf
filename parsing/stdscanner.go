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

// The standard scanner. It uses the "text/scanner" package.
type StdScanner struct{
	s    scanner.Scanner
	tmap map[rune]string
	kmap map[string]string
}
/*
 Initializes the Scanner.

 kmap map[string]string:  the keyword map. Initialize with SyntaxFile.ScanForKeyWords

 src io.Reader:           the source code.
 */
func (s *StdScanner) Init(kmap map[string]string,src io.Reader) *StdScanner {
	s.tmap = stdTMap
	s.kmap = kmap
	if s.kmap==nil { s.kmap = map[string]string{} }
	s.s.Init(src)
	return s
}

/*
 Returns the first token. This method should only called once.
 The names of the token types are derived from the "text/scanner" package.
 */
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


