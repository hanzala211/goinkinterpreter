package scanner

import (
	"strconv"

	"github.com/hanzala211/goinkinterpreter/token"
)

type Scanner struct {
	source []rune
	Tokens []*token.Token
	vm     vm

	current int
	start   int
	line    int
}

type vm interface {
	ReportError(line int, err error)
}

func NewScanner(source string, vm vm) *Scanner {
	return &Scanner{
		source: []rune(source),
		vm:     vm,
		line:   1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, &token.Token{
		Type:    token.TokenType_EOF,
		Line:    s.line,
		Lexeme:  "",
		Literal: nil,
	})
}

func (s *Scanner) scanToken() {
	var c rune = s.advance()
	switch c {
	case ' ', '\t', '\r':
	case '\n':
		s.line++
	case '+':
		s.addToken(token.TokenType_Plus)
	case '-':
		s.addToken(token.TokenType_Minus)
	case '(':
		s.addToken(token.TokenType_LeftParen)
	case ')':
		s.addToken(token.TokenType_RightParen)
	case '{':
		s.addToken(token.TokenType_LeftBrace)
	case '}':
		s.addToken(token.TokenType_RightBrace)
	case ',':
		s.addToken(token.TokenType_Comma)
	case '.':
		s.addToken(token.TokenType_Dot)
	case ';':
		s.addToken(token.TokenType_Semicolon)
	case '*':
		s.addToken(token.TokenType_Star)

	case '!':
		if s.match('=') {
			s.addToken(token.TokenType_BangEqual)
		} else {
			s.addToken(token.TokenType_Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.TokenType_EqualEqual)
		} else {
			s.addToken(token.TokenType_Equal)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.TokenType_GreaterEqual)
		} else {
			s.addToken(token.TokenType_Greater)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.TokenType_LessEqual)
		} else {
			s.addToken(token.TokenType_Less)
		}
	case '/':
		if s.match('/') {
			s.consumeLineComment()
			break
		}
		if s.match('*') {
			s.consumeBlockComment()
			break
		}
		s.addToken(token.TokenType_Slash)
	case '"':
		s.scanString()

	default:
		if s.isDigit(c) {
			s.scanNumber()
		} else if s.isAlpha(c) {
			s.scanIdentifier()
		} else {
			return
		}
	}
}
func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal any) {
	s.Tokens = append(s.Tokens, &token.Token{
		Type:    tokenType,
		Line:    s.line,
		Lexeme:  string(s.source[s.start:s.current]),
		Literal: literal,
	})
}

func (s *Scanner) match(c rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != c {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) consumeLineComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) consumeBlockComment() {
	for !s.isAtEnd() {
		switch c := s.advance(); c {
		case '*':
			if s.peekNext() == '/' {
				s.current += 2
				return
			}
		case '\n':
			s.line++
			s.advance()
		default:
			s.advance()
		}
	}
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return
	}

	s.advance()
	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(token.TokenType_String, value)
}
func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) scanNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value := string(s.source[s.start:s.current])
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.addTokenWithLiteral(token.TokenType_Number, num)
}

func (s *Scanner) scanIdentifier() {
	for s.isAlpha(s.peek()) || s.isDigit(s.peek()) {
		s.advance()
	}
	text := string(s.source[s.start:s.current])
	if t, ok := keywords[text]; ok {
		s.addToken(t)
		return
	}

	s.addToken(token.TokenType_Identifier)
}
