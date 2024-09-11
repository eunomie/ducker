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
	case token.ARG:
		return p.parseArgStatement()
	case token.COPY:
		return p.parseCopyStatement()
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

	p.expectPeekEnd()

	return stmt
}

func (p *Parser) parseArguments() []*ast.Argument {
	var args []*ast.Argument

	for strings.HasPrefix(p.curToken.Literal, "--") {
		a := strings.TrimPrefix(p.curToken.Literal, "--")
		if name, value, ok := strings.Cut(a, "="); ok {
			args = append(args, &ast.Argument{Token: p.curToken, Name: name, Value: value})
		} else {
			return nil
		}
		if !p.expectPeek(token.EXPR) {
			return nil
		}
	}

	return args
}

func (p *Parser) parseFromStatement() *ast.FromStatement {
	stmt := &ast.FromStatement{Token: p.curToken}

	if !p.expectPeek(token.EXPR) {
		return nil
	}

	stmt.Arguments = p.parseArguments()

	stmt.Source = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.AS) {
		p.nextToken()
		p.nextToken()

		stmt.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	p.expectPeekEnd()

	return stmt
}

func (p *Parser) parseCopyStatement() *ast.CopyStatement {
	stmt := &ast.CopyStatement{Token: p.curToken}

	if !p.expectPeek(token.EXPR) {
		return nil
	}

	stmt.Arguments = p.parseArguments()

	stmt.Source = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.EXPR) {
		return nil
	}

	stmt.Dest = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.expectPeekEnd()

	return stmt
}

func (p *Parser) parseArgStatement() *ast.ArgStatement {
	stmt := &ast.ArgStatement{Token: p.curToken}

	if !p.expectPeek(token.EXPR) {
		return nil
	}

	stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if name, value, ok := strings.Cut(p.curToken.Literal, "="); ok {
		stmt.ArgName = name
		stmt.ArgValue = strPtr(value)
	} else {
		stmt.ArgName = name
		if p.peekTokenIs(token.EXPR) {
			p.nextToken()
			stmt.ArgValue = strPtr(p.curToken.Literal)
		}
	}

	p.expectPeekEnd()

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

func (p *Parser) expectPeekEnd() {
	if !p.peekTokenIs(token.EOL) && !p.peekTokenIs(token.EOF) {
		p.peekError(token.EOL)
	} else {
		p.nextToken()
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func strPtr(s string) *string {
	return &s
}
