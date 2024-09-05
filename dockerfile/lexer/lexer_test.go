package lexer

import (
	"fmt"
	"testing"

	"github.com/eunomie/ducker/dockerfile/token"

	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	input := `
# syntax=docker/dockerfile:1
# check=error=true

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx
FROM alpine AS builder
ENV GOPRIVATE=github.com/eunomie
COPY --from=xx / /
RUN apk add --no-cache curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    task go:build && xx-verify dist/scout-notifications-service

COPY --from=build-base /etc/ssl/certs/ca\ certificates.crt /etc/ssl/certs/

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/scout-notifications-service"]

CMD ["echo", "Hello, World!"]
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.EOL, "\n"},
		{token.COMMENT, "#"},
		{token.EXPR, "syntax=docker/dockerfile:1"},
		{token.EOL, "\n"},
		{token.COMMENT, "#"},
		{token.EXPR, "check=error=true"},
		{token.EOL, "\n"},
		{token.EOL, "\n"},
		{token.FROM, "FROM"},
		{token.EXPR, "--platform=$BUILDPLATFORM"},
		{token.EXPR, "tonistiigi/xx:${XX_VERSION}"},
		{token.AS, "AS"},
		{token.EXPR, "xx"},
		{token.EOL, "\n"},
		{token.FROM, "FROM"},
		{token.EXPR, "alpine"},
		{token.AS, "AS"},
		{token.EXPR, "builder"},
		{token.EOL, "\n"},
		{token.ENV, "ENV"},
		{token.EXPR, "GOPRIVATE=github.com/eunomie"},
		{token.EOL, "\n"},
		{token.COPY, "COPY"},
		{token.EXPR, "--from=xx"},
		{token.EXPR, "/"},
		{token.EXPR, "/"},
		{token.EOL, "\n"},
		{token.RUN, "RUN"},
		{token.EXPR, "apk"},
		{token.EXPR, "add"},
		{token.EXPR, "--no-cache"},
		{token.EXPR, "curl"},
		{token.EOL, "\n"},
		{token.RUN, "RUN"},
		{token.EXPR, "sh"},
		{token.EXPR, "-c"},
		{token.STRING, "$(curl --location https://taskfile.dev/install.sh)"},
		{token.EXPR, "--"},
		{token.EXPR, "-d"},
		{token.EXPR, "-b"},
		{token.EXPR, "/usr/local/bin"},
		{token.EOL, "\n"},
		{token.RUN, "RUN"},
		{token.EXPR, "--mount=type=cache,target=/root/.cache"},
		{token.EXPR, "--mount=type=cache,target=/go/pkg/mod"},
		{token.EXPR, "task"},
		{token.EXPR, "go:build"},
		{token.EXPR, "&&"},
		{token.EXPR, "xx-verify"},
		{token.EXPR, "dist/scout-notifications-service"},
		{token.EOL, "\n"},
		{token.EOL, "\n"},
		{token.COPY, "COPY"},
		{token.EXPR, "--from=build-base"},
		{token.EXPR, "/etc/ssl/certs/ca\\ certificates.crt"},
		{token.EXPR, "/etc/ssl/certs/"},
		{token.EOL, "\n"},
		{token.EOL, "\n"},
		{token.EXPOSE, "EXPOSE"},
		{token.EXPR, "8080"},
		{token.EOL, "\n"},
		{token.ENTRYPOINT, "ENTRYPOINT"},
		{token.LBRACKET, "["},
		{token.STRING, "/usr/local/bin/scout-notifications-service"},
		{token.RBRACKET, "]"},
		{token.EOL, "\n"},
		{token.EOL, "\n"},
		{token.CMD, "CMD"},
		{token.LBRACKET, "["},
		{token.STRING, "echo"},
		{token.COMMA, ","},
		{token.STRING, "Hello, World!"},
		{token.RBRACKET, "]"},
		{token.EOL, "\n"},
	}

	l := New(input)

	for i, test := range tests {
		tt := test
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tok := l.NextToken()

			require.Equal(t, tt.expectedType, tok.Type)
			require.Equal(t, tt.expectedLiteral, tok.Literal)
		})
	}
}
