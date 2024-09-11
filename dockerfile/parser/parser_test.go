package parser

import (
	"fmt"
	"testing"

	"github.com/eunomie/ducker/dockerfile/ast"
	"github.com/eunomie/ducker/dockerfile/lexer"
	"github.com/stretchr/testify/require"
)

func TestFromStatements(t *testing.T) {
	input := `
FROM alpine
FROM docker.io/library/golang:1.23.0 AS builder
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)

	tests := []struct {
		expectedSource string
		expectedAlias  *string
		expectedArgs   []string
	}{
		{"alpine", nil, nil},
		{"docker.io/library/golang:1.23.0", strPtr("builder"), nil},
		{"tonistiigi/xx:${XX_VERSION}", strPtr("xx"), []string{"--platform=$BUILDPLATFORM"}},
	}

	for i, test := range tests {
		tt := test
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			stmt := program.Statements[i]
			testFromStatement(t, stmt, tt.expectedSource, tt.expectedAlias, tt.expectedArgs)
		})
	}
}

func testFromStatement(t *testing.T, s ast.Statement, expectedSource string, expectedAlias *string, expectedArgs []string) {
	t.Helper()

	require.Equal(t, "FROM", s.TokenLiteral())
	require.IsType(t, &ast.FromStatement{}, s)

	fromStmt, _ := s.(*ast.FromStatement)
	require.Equal(t, expectedSource, fromStmt.Source.Value)
	require.Equal(t, expectedSource, fromStmt.Source.TokenLiteral())

	if expectedAlias != nil {
		require.NotNil(t, fromStmt.Alias)
		require.Equal(t, *expectedAlias, fromStmt.Alias.Value)
		require.Equal(t, *expectedAlias, fromStmt.Alias.TokenLiteral())
	} else {
		require.Nil(t, fromStmt.Alias)
	}

	var args []string
	for _, arg := range fromStmt.Arguments {
		args = append(args, arg.TokenLiteral())
	}
	require.ElementsMatch(t, expectedArgs, args)
}

func strPtr(s string) *string {
	return &s
}
