package ast

import "github.com/maxymania/langtransf/parsing"
import "fmt"

type AST struct{
	Head string
	Content string
	Token *parsing.Token
	Childs []*AST
}
func (a *AST)NewAst(head,content string,token *parsing.Token) *AST{
	if a==nil { return nil }
	return &AST{head,content,token,nil}
}
func (a *AST)Add(child *AST){
	if a==nil { return }
	a.Childs = append(a.Childs,child)
}

func (a *AST)Backup() (b AST){
	if a==nil { return }
	b = *a
	return
}
func (a *AST)Restore(b AST){
	if a==nil { return }
	*a = b
}
func (a *AST)String() string{
	if a==nil { return "nil" }
	if len(a.Content)!=0 { return fmt.Sprint("\"",a.Content,"\"") }
	return fmt.Sprint("~`",a.Head,"`",a.Childs)
}


