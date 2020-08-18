package parser

/*
   Creation Time: 2020 - Aug - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// token identifies the type of lex items.
type token int

const (
	ERROR token = iota
	EOF
	WHITE_SPACE

	literal_beg
	IDENT // names, variables
	literal_end

	operator_beg
	AT_SIGN // @
	COMMA   // ,
	LPAREN  // (
	RPAREN  // )
	operator_end

	keyword_beg
	MODEL
	TABLE
	VIEW
	COUNTER
	ATOMIC_COUNTER
	keyword_end
)

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }
