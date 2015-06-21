package parsing

import "sync"

type Generator func() *Token

/*
 Represents a token resolved from a scanner.
 */
type Token struct{
	t,d,txt string
	pos interface{}
	initf func()
	inits sync.Once
	nxt  *Token
}

func noop() { }

func initialize(gen Generator, t *Token) func() {
	return func(){
		t.nxt = gen()
	}
}

func MakeToken(
		typ,distinct,text string,
		pos interface{},nxt Generator ) (t *Token) {
	t = new(Token)
	t.t   = typ
	t.d   = distinct
	t.txt = text
	if nxt==nil {
		t.initf = noop
	} else {
		t.initf = initialize(nxt,t)
	}
	return
}

/*
 Returns the general token type, if the token is not a keyword and not a
 distinct one-character-token.
 */
func (tl *Token) Type() string { return tl.t }

/*
 Returns the token text, if the token is a keyword or a
 distinct one-character-token.
 */
func (tl *Token) Distinct() string { return tl.d }

/*
 Returns the token text.
 */
func (tl *Token) Text() string { return tl.txt }

/*
 Returns the token text position.
 */
func (tl *Token) Position() interface{} { return tl.pos }

/*
 Returns the next token following to this.
 */
func (tl *Token) Next() *Token {
	tl.inits.Do(tl.initf)
	return tl.nxt
}


