//go:generate stringer -type=TokenType -trimprefix=TokenType_

package token

type TokenType int

const (
	// Single-character tokens.
	TokenType_LeftParen  TokenType = iota // (
	TokenType_RightParen                  // )
	TokenType_LeftBrace                   // {
	TokenType_RightBrace                  // }
	TokenType_Comma                       // ,
	TokenType_Dot                         // .
	TokenType_Minus                       // -
	TokenType_Plus                        // +
	TokenType_Semicolon                   // ;
	TokenType_Slash                       // /
	TokenType_Star                        // *

	// One or two character tokens.
	TokenType_Bang         // !
	TokenType_BangEqual    // !=
	TokenType_Equal        // =
	TokenType_EqualEqual   // ==
	TokenType_Greater      // >
	TokenType_GreaterEqual // >=
	TokenType_Less         // <
	TokenType_LessEqual    // <=

	// Literals.
	TokenType_Identifier // a
	TokenType_String     // "a"
	TokenType_Number     // 1

	// Keywords.
	TokenType_And    // and
	TokenType_Class  // class
	TokenType_Else   // else
	TokenType_False  // false
	TokenType_Fun    // fun
	TokenType_For    // for
	TokenType_If     // if
	TokenType_Nil    // nil
	TokenType_Or     // or
	TokenType_Print  // print
	TokenType_Return // return
	TokenType_Super  // super
	TokenType_This   // this
	TokenType_True   // true
	TokenType_Var    // var
	TokenType_While  // while

	TokenType_EOF // eof
)
