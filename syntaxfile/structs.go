package syntaxfile

import "fmt"
import "bytes"
import "github.com/maxymania/langtransf/parsing"
import "github.com/maxymania/langtransf/ast"

const (
	M_OR byte = iota
	M_SEQ
	M_DROP
	M_OMIT_VERBOSITY
)
const (
	F_VERBOSE byte = 1<<iota
	F_MUTE
)

var m_syms = []string{
	"or",
	"",
	"#-",
}

func generateSRule(r []Rule) string {
	b := &bytes.Buffer{}
	b.WriteString("(")
	for i,v := range r {
		if i==0{
			fmt.Fprintf(b,"%v",v)
		}else{
			fmt.Fprintf(b," %v",v)
		}
	}
	b.WriteString(")")
	return b.String()
}

func generateSRule_or(r []Rule) string {
	b := &bytes.Buffer{}
	b.WriteString("(")
	for i,v := range r {
		if i==0{
			fmt.Fprintf(b,"%v",v)
		}else{
			fmt.Fprintf(b,"|%v",v)
		}
	}
	b.WriteString(")")
	return b.String()
}

type Rule interface{
	ScanForKeyWords(km map[string]string)
	Parse(sf SyntaxFile,t *parsing.Token,d *ast.AST,e *ErrorRecorder) *parsing.Token 
}

type SyntaxFile map[string]Rule
func (sf SyntaxFile) ScanForKeyWords(km map[string]string) {
	for _,r := range sf { r.ScanForKeyWords(km) }
}
func (sf SyntaxFile) String() string {
	b := &bytes.Buffer{}
	for k,v := range sf {
		fmt.Fprintf(b,"%s = %v ;\n",k,v)
	}
	return b.String()
}

type Modifier struct{
	Data []Rule
	Mode byte
	Flags byte
}
func (m Modifier) String() string {
	if m.Mode==M_OR {
		return fmt.Sprintf("%v",generateSRule_or(m.Data))
	}
	if m.Mode==M_SEQ {
		return fmt.Sprintf("%v",generateSRule(m.Data))
	}
	if m.Mode==M_DROP {
		return fmt.Sprintf("%v %v",generateSRule(m.Data),m_syms[m.Mode])
	}
	if m.Mode==M_OMIT_VERBOSITY	{
		return fmt.Sprintf("%v #p",m.Data[0])
	}
	if m.Mode<' ' {
		return fmt.Sprintf("%v",generateSRule(m.Data))
	}
	return fmt.Sprintf("%v%c",generateSRule(m.Data),m.Mode)
}
func (m Modifier) ScanForKeyWords(km map[string]string) {
	for _,r := range m.Data { r.ScanForKeyWords(km) }
}
func (m Modifier) Parse(sf SyntaxFile,t *parsing.Token, d *ast.AST,e *ErrorRecorder) *parsing.Token {
	if (m.Flags&F_MUTE)!=0 { e=nil }
	switch (m.Mode){
	case M_OMIT_VERBOSITY:
		{
			ebak := e.Backup()
			t = m.Data[0].Parse(sf,t,d,e)
			if t!=nil { e.Restore(ebak) }
			return t
		}
	case M_DROP:
		d = nil
		fallthrough
	case M_SEQ:
		for _,r := range m.Data {
			t = r.Parse(sf,t,d,e)
			if t==nil { return nil }
		}
		return t
	case M_OR:
		ebak := e.Backup()
		for _,r := range m.Data {
			if (m.Flags&F_VERBOSE)==0 { e.Restore(ebak) }
			bak := d.Backup()
			ebak = e.Backup()
			ret := r.Parse(sf,t,d,e)
			if ret!=nil { return ret }
			d.Restore(bak)
		}
	case '?':
		{
			bak := d.Backup()
			ebak := e.Backup()
			ret := m.Data[0].Parse(sf,t,d,e)
			if ret!=nil { return ret }
			d.Restore(bak)
			if (m.Flags&F_VERBOSE)==0 { e.Restore(ebak) }
			return t
		}
	case '+':
		t = m.Data[0].Parse(sf,t,d,e)
		if t==nil { return nil }
		fallthrough
	case '*':
		for {
			bak := d.Backup()
			ebak := e.Backup()
			ret := m.Data[0].Parse(sf,t,d,e)
			if ret!=nil { t = ret; continue }
			d.Restore(bak)
			if (m.Flags&F_VERBOSE)==0 { e.Restore(ebak) }
			return t
		}
	}
	return nil
}

type Group struct{
	Data Rule
	Name string
}
func (g Group) String() string {
	return fmt.Sprintf("~(`%v` %v )",g.Name,g.Data)
}
func (g Group) ScanForKeyWords(km map[string]string) {
	g.Data.ScanForKeyWords(km)
}
func (g Group) Parse(sf SyntaxFile,t *parsing.Token, d *ast.AST,e *ErrorRecorder) *parsing.Token {
	d2 := d.NewAst(g.Name,"",nil)
	r := g.Data.Parse(sf,t,d2,e)
	if r==nil { return nil }
	d.Add(d2)
	return r
}

type MatchText struct{
	Text string
}
func (m MatchText) String() string {
	return fmt.Sprintf("\"%s\"",m.Text)
}
func (m MatchText) ScanForKeyWords(km map[string]string) {
	km[m.Text]=m.Text
}
func (m MatchText) Parse(sf SyntaxFile,t *parsing.Token, d *ast.AST,e *ErrorRecorder) *parsing.Token {
	if t.Distinct()!=m.Text {
		if t.Distinct()!="" {
			e.AddErr(t.Position(),
				fmt.Sprintf("expected \"%s\", got \"%s\"",m.Text,t.Distinct()))
		}else{
			e.AddErr(t.Position(),
				fmt.Sprintf("expected \"%s\", got %s",m.Text,t.Type()))
		}
		return nil
	}
	d.Add(d.NewAst("",t.Text(),t))
	return t.Next()
}

type MatchToken struct{
	Token string
}
func (m MatchToken) String() string {
	return m.Token
}
func (m MatchToken) ScanForKeyWords(km map[string]string) {}
func (m MatchToken) Parse(sf SyntaxFile,t *parsing.Token, d *ast.AST,e *ErrorRecorder) *parsing.Token {
	if t.Type()!=m.Token {
		if t.Distinct()!="" {
			e.AddErr(t.Position(),
				fmt.Sprintf("expected %s, got \"%s\"",m.Token,t.Distinct()))
		}else{
			e.AddErr(t.Position(),
				fmt.Sprintf("expected %s, got %s",m.Token,t.Type()))
		}
		return nil
	}
	d.Add(d.NewAst("",t.Text(),t))
	return t.Next()
}

type CallRule struct{
	Name string
}
func (c CallRule) String() string {
	return c.Name
}
func (c CallRule) ScanForKeyWords(km map[string]string) {}
func (c CallRule) Parse(sf SyntaxFile,t *parsing.Token, d *ast.AST,e *ErrorRecorder) *parsing.Token {
	r,ok := sf[c.Name]
	if !ok { panic("error: rune not found: "+c.Name) }
	return r.Parse(sf,t,d,e)
}

