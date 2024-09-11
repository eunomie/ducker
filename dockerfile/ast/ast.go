package ast

import "github.com/eunomie/ducker/dockerfile/token"

type (
	Node interface {
		TokenLiteral() string
	}

	Statement interface {
		Node
		statementNode()
	}

	Expression interface {
		Node
		expressionNode()
	}
)

type Program struct {
	Directive  []Statement
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type CommentStatement struct {
	Token token.Token
	Value *Identifier
}

func (fs *CommentStatement) statementNode()       {}
func (fs *CommentStatement) TokenLiteral() string { return fs.Token.Literal }

type FromStatement struct {
	Token     token.Token
	Source    *Identifier
	Alias     *Identifier
	Arguments []*Argument
}

func (fs *FromStatement) statementNode()       {}
func (fs *FromStatement) TokenLiteral() string { return fs.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type Argument struct {
	Token token.Token
	Name  string
	Value string
}

func (a *Argument) expressionNode()      {}
func (a *Argument) TokenLiteral() string { return a.Token.Literal }

type ArgStatement struct {
	Token    token.Token
	Value    *Identifier
	ArgName  string
	ArgValue *string
}

func (fs *ArgStatement) statementNode()       {}
func (fs *ArgStatement) TokenLiteral() string { return fs.Token.Literal }

type CopyStatement struct {
	Token     token.Token
	Source    *Identifier
	Dest      *Identifier
	Arguments []*Argument
}

func (fs *CopyStatement) statementNode()       {}
func (fs *CopyStatement) TokenLiteral() string { return fs.Token.Literal }
