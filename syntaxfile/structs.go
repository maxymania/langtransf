package syntaxfile

import "fmt"
import "bytes"

const (
	M_OR byte = iota
	M_SEQ
	M_DROP
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

type Rule interface{}

type SyntaxFile map[string]Rule
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
}
func (m Modifier) String() string {
	if m.Mode==0 {
		return fmt.Sprintf("%v",generateSRule_or(m.Data))
	}
	if m.Mode==1 {
		return fmt.Sprintf("%v",generateSRule(m.Data))
	}
	if m.Mode==2 {
		return fmt.Sprintf("%v %v",generateSRule(m.Data),m_syms[m.Mode])
	}
	if m.Mode<' ' {
		return fmt.Sprintf("%v",generateSRule(m.Data))
	}
	return fmt.Sprintf("%v%c",generateSRule(m.Data),m.Mode)
}

type Group struct{
	Data Rule
	Name string
}
func (g Group) String() string {
	return fmt.Sprintf("~(`%v` %v )",g.Name,g.Data)
}

type MatchText struct{
	Text string
}
func (m MatchText) String() string {
	return fmt.Sprintf("\"%d\"",m.Text)
}

type MatchToken struct{
	Token string
}
func (m MatchToken) String() string {
	return m.Token
}

type CallRule struct{
	Name string
}
func (c CallRule) String() string {
	return c.Name
}

