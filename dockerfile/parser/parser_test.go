package parser

import (
	"fmt"
	"testing"

	"github.com/eunomie/ducker/dockerfile/ast"
	"github.com/eunomie/ducker/dockerfile/lexer"
	"github.com/stretchr/testify/require"
)

func TestDockerfile(t *testing.T) {
	input := `# syntax=docker/dockerfile:1
# check=error=true

ARG XX_VERSION=1.2.1

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:1.22.5-alpine3.20 AS build-base
COPY --from=xx / /
`

	program := requireParse(t, input)
	require.Len(t, program.Directive, 2)
	require.Len(t, program.Statements, 4)
}

func TestArgStatement(t *testing.T) {
	input := `
ARG VERSION=latest

ARG BUILDPLATFORM

ARG XX_VERSION 1.2.1`

	program := requireParse(t, input)
	require.Len(t, program.Statements, 3)

	tests := []struct {
		expectedName  string
		expectedValue *string
	}{
		{"VERSION", strPtr("latest")},
		{"BUILDPLATFORM", nil},
		{"XX_VERSION", strPtr("1.2.1")},
	}

	for i, test := range tests {
		tt := test
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			stmt := program.Statements[i]
			testArgStatement(t, stmt, tt.expectedName, tt.expectedValue)
		})
	}
}

func testArgStatement(t *testing.T, s ast.Statement, expectedName string, expectedValue *string) {
	t.Helper()

	require.Equal(t, "ARG", s.TokenLiteral())
	require.IsType(t, &ast.ArgStatement{}, s)

	argStmt, _ := s.(*ast.ArgStatement)
	require.Equal(t, expectedName, argStmt.ArgName)

	if expectedValue != nil {
		require.NotNil(t, argStmt.ArgValue)
		require.Equal(t, *expectedValue, *argStmt.ArgValue)
	} else {
		require.Nil(t, argStmt.ArgValue)
	}
}

func TestFromStatements(t *testing.T) {
	input := `
FROM alpine
FROM docker.io/library/golang:1.23.0 AS builder
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx`

	program := requireParse(t, input)
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

func TestCopyStatement(t *testing.T) {
	input := `
COPY --from=xx / /
COPY --from=build-base /etc/ssl/certs/ca\ certificates.crt /etc/ssl/certs/`

	program := requireParse(t, input)
	require.Len(t, program.Statements, 2)

	tests := []struct {
		expectedSource string
		expectedTarget string
		expectedArgs   []string
	}{
		{"/", "/", []string{"--from=xx"}},
		{"/etc/ssl/certs/ca\\ certificates.crt", "/etc/ssl/certs/", []string{"--from=build-base"}},
	}

	for i, test := range tests {
		tt := test
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			stmt := program.Statements[i]
			testCopyStatement(t, stmt, tt.expectedSource, tt.expectedTarget, tt.expectedArgs)
		})
	}
}

func testCopyStatement(t *testing.T, s ast.Statement, expectedSource, expectedTarget string, expectedArgs []string) {
	t.Helper()

	require.Equal(t, "COPY", s.TokenLiteral())
	require.IsType(t, &ast.CopyStatement{}, s)

	copyStmt, _ := s.(*ast.CopyStatement)
	require.Equal(t, expectedSource, copyStmt.Source.Value)
	require.Equal(t, expectedSource, copyStmt.Source.TokenLiteral())
	require.Equal(t, expectedTarget, copyStmt.Dest.Value)
	require.Equal(t, expectedTarget, copyStmt.Dest.TokenLiteral())

	var args []string
	for _, arg := range copyStmt.Arguments {
		args = append(args, arg.TokenLiteral())
	}
	require.ElementsMatch(t, expectedArgs, args)
}

func requireParse(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	return program
}
