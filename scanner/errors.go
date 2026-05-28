package scanner

import "errors"

var (
	ErrUnexpectedEOF          = errors.New("unexpected EOF")
	ErrUnexpectedRune         = errors.New("unexpected rune")
	ErrUnindentedBlockComment = errors.New("unindented block comment")
	ErrUnterminatedComment    = errors.New("Multiline comment met EOF.")
	ErrUnterminatedString     = errors.New("Unterminated string.")
	ErrInvalidNumber          = errors.New("Invalid number.")
)

