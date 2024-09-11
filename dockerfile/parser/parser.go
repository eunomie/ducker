package parser

import (
	"fmt"
	"strings"

	"github.com/eunomie/ducker/dockerfile/ast"
	"github.com/eunomie/ducker/dockerfile/lexer"
	"github.com/eunomie/ducker/dockerfile/token"
)

type (
	Parser struct {
		l *lexer.Lexer

		errors []string

		curToken  token.Token
		peekToken token.Token
	}
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			if s, ok := stmt.(*ast.CommentStatement); ok && len(program.Statements) == 0 {
				program.Directive = append(program.Directive, s)
			} else {
				program.Statements = append(program.Statements, stmt)
			}
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.COMMENT:
		return p.parseCommentStatement()
	case token.FROM:
		return p.parseFromStatement()
	default:
		return nil
	}
}

func (p *Parser) parseCommentStatement() *ast.CommentStatement {
	stmt := &ast.CommentStatement{Token: p.curToken}

	if !p.expectPeek(token.EXPR) {
		return nil
	}

	stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

func (p *Parser) parseFromStatement() *ast.FromStatement {
	stmt := &ast.FromStatement{Token: p.curToken}

	if !p.expectPeek(token.EXPR) {
		return nil
	}

	if strings.HasPrefix(p.curToken.Literal, "--") {
		a := strings.TrimPrefix(p.curToken.Literal, "--")
		if name, value, ok := strings.Cut(a, "="); ok {
			stmt.Arguments = append(stmt.Arguments, &ast.Argument{Token: p.curToken, Name: name, Value: value})
		} else {
			return nil
		}
		if !p.expectPeek(token.EXPR) {
			return nil
		}
	}

	stmt.Source = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.AS) {
		p.nextToken()
		p.nextToken()

		stmt.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	//...

	for !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
