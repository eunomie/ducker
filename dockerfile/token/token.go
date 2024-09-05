package token

type (
	Type string

	Token struct {
		Type     Type
		Literal  string
		Position Position
	}

	Position struct {
		Line  int
		Start int
	}
)

const (
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"

	EOL     Type = "EOL"
	EXPR    Type = "EXPR"
	STRING  Type = "STRING"
	COMMENT Type = "COMMENT"

	LBRACKET Type = "["
	RBRACKET Type = "]"
	COMMA    Type = ","

	// all the keywords

	FROM       Type = "FROM"
	AS         Type = "AS"
	USER       Type = "USER"
	ENV        Type = "ENV"
	ARG        Type = "ARG"
	COPY       Type = "COPY"
	RUN        Type = "RUN"
	WORKDIR    Type = "WORKDIR"
	EXPOSE     Type = "EXPOSE"
	ENTRYPOINT Type = "ENTRYPOINT"
	CMD        Type = "CMD"
)

var (
	keywords = map[string]Type{
		"FROM":       FROM,
		"AS":         AS,
		"USER":       USER,
		"ENV":        ENV,
		"ARG":        ARG,
		"COPY":       COPY,
		"RUN":        RUN,
		"WORKDIR":    WORKDIR,
		"EXPOSE":     EXPOSE,
		"ENTRYPOINT": ENTRYPOINT,
		"CMD":        CMD,
	}
)

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return EXPR
}
