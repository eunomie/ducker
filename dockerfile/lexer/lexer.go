package lexer

import (
	"unicode/utf8"

	"github.com/eunomie/ducker/dockerfile/token"
)

type (
	Lexer struct {
		input        string
		position     int
		readPosition int
		line         int
		lineStart    int
		r            rune
	}
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readRune()
	return l
}

func (l *Lexer) readRune() {
	size := 1
	if l.readPosition >= len(l.input) {
		l.r = utf8.RuneError
	} else {
		l.r, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.position = l.readPosition
	l.readPosition += size
}

func (l *Lexer) PeekRune() rune {
	if l.readPosition >= len(l.input) {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.r {
	case '#':
		tok = newToken(token.COMMENT, l.r)
	case '[':
		tok = newToken(token.LBRACKET, l.r)
	case ']':
		tok = newToken(token.RBRACKET, l.r)
	case '"':
		tok = newToken(token.STRING, l.r)
		tok.Literal = l.readString('"')
	case '\'':
		tok = newToken(token.STRING, l.r)
		tok.Literal = l.readString('\'')
	case '\n':
		tok = newToken(token.EOL, l.r)
		l.line += 1
		l.lineStart = l.readPosition
	case ',':
		tok = newToken(token.COMMA, l.r)
	case utf8.RuneError:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if l.isValidExprRune() {
			tok.Position = l.userPosition()
			tok.Literal = l.readExpression()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.r)
		}
	}

	tok.Position = l.userPosition()

	l.readRune()
	return tok
}

func (l *Lexer) readExpression() string {
	position := l.position
	for l.isValidExprRune() {
		l.readRune()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(del rune) string {
	l.readRune()
	position := l.position
	for l.r != del {
		l.readRune()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readFlag() string {
	l.readRune()
	position := l.position
	for l.r != ' ' && l.r != '\t' && l.r != '\n' && l.r != '\r' && l.r != utf8.RuneError {
		l.readRune()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.r {
		case ' ', '\t', '\r':
		case '\\':
			if l.PeekRune() == '\n' {
				l.readRune()
				l.line += 1
				l.lineStart = l.readPosition
			}
		default:
			return
		}
		l.readRune()
	}
}

func (l *Lexer) userPosition() token.Position {
	return token.Position{Line: l.line, Start: l.position - l.lineStart}
}

func (l *Lexer) isValidExprRune() bool {
	switch l.r {
	case ' ', '\t', '\r', '\n', utf8.RuneError:
		return false
	case '\\':
		if l.PeekRune() == ' ' {
			l.readRune()
			return true
		}
		return false
	default:
		return true
	}
}

func newToken(tokenType token.Type, r rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(r)}
}
