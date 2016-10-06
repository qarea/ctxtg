//Package ctxtgtest provides test utilities for ctgtg users
package ctxtgtest

import (
	"errors"
	"time"

	"github.com/qarea/ctxtg"
)

// Check if it implements Parser/Signer interfaces
var (
	_ ctxtg.TokenParser = (*Parser)(nil)
	_ ctxtg.TokenSigner = (*Signer)(nil)
)

// Errors for ctxtgtest package
var (
	ErrMethodNotCalled   = errors.New("Method wasn't called as expected")
	ErrUnexpectedToken   = errors.New("Unexpected token passed")
	ErrUnexpectedClaims  = errors.New("Unexpected Claims passed")
	ErrUnexpectedTimeout = errors.New("Unexpected timeout passed")
)

// Parser implements ctxtg.TokenParser interface for parser mocking for Unit tests
type Parser struct {
	//Values to return on Parse call
	Claims        ctxtg.Claims
	Err           error
	TokenExpected ctxtg.Token

	token  ctxtg.Token
	called bool
}

// ParseCtxWithClaims use Parser.Parse function under the hood
func (p *Parser) ParseCtxWithClaims(context ctxtg.Context, f ctxtg.CtxClaimsFunc) error {
	c, err := p.Parse(context.Token)
	if err != nil {
		return err
	}
	ctx, cancel := context.ToContext()
	defer cancel()
	return f(ctx, *c)
}

// ParseWithClaims use Parser.Parse function under the hood
func (p *Parser) ParseWithClaims(token ctxtg.Token, f ctxtg.ClaimsFunc) error {
	c, err := p.Parse(token)
	if err != nil {
		return err
	}
	return f(*c)
}

// Parse save token to Parser and register fact of calling this method
func (p *Parser) Parse(token ctxtg.Token) (*ctxtg.Claims, error) {
	p.called = true
	p.token = token
	return &p.Claims, p.Err
}

// Error return error if Parse method wasn't called or token passed to Parse wasn't expected
func (p *Parser) Error() error {
	if !p.called {
		return ErrMethodNotCalled
	}
	if p.TokenExpected != p.token {
		return ErrUnexpectedToken
	}
	return nil
}

// Signer implements ctxtg.TokenSigner interface for signer mocking for Unit tests
type Signer struct {
	// Values to return on Sign call
	Token ctxtg.Token
	Err   error

	ClaimsExpected  ctxtg.Claims
	TimeoutExpected time.Duration

	timeout time.Duration
	claims  ctxtg.Claims
	called  bool
}

// Sign save args to Signer, register call and return values from Token and Err fields
func (s *Signer) Sign(c ctxtg.Claims, timeout time.Duration) (ctxtg.Token, error) {
	s.called = true
	s.claims = c
	s.timeout = timeout
	return s.Token, s.Err
}

// Error returns err if Sign method wasn't called or unexpected arguments passed to Sign wasn't expected
func (s *Signer) Error() error {
	if !s.called {
		return ErrMethodNotCalled
	}
	if s.ClaimsExpected != s.claims {
		return ErrUnexpectedClaims
	}
	if s.TimeoutExpected != s.timeout {
		return ErrUnexpectedTimeout
	}
	return nil
}
