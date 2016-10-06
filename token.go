package ctxtg

import (
	"crypto/rsa"
	"strconv"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Global errors for all projects
var (
	ErrInvalidToken = jsonrpc2.NewError(1, "INVALID_TOKEN")
	ErrTokenExpired = jsonrpc2.NewError(2, "TOKEN_EXPIRED")
)

var timeNowFunc = time.Now

// Claims represents encoded into JWT info
type Claims struct {
	UserID UserID
}

// UserID represents user id in Timeguard system
type UserID int64

// ClaimsFunc is function in which claims will be passed if JWT Token is fine
type ClaimsFunc func(Claims) error

// CtxClaimsFunc is function in which claims and converted context.Context will be passed if JWT Token is fine
type CtxClaimsFunc func(context.Context, Claims) error

// TokenParser interface for JWT token parsers and point for mocking (see ctxtgtest subpackage)
type TokenParser interface {
	Parse(Token) (*Claims, error)
	ParseWithClaims(Token, ClaimsFunc) error
	ParseCtxWithClaims(Context, CtxClaimsFunc) error
}

// NewRSATokenParser parse publicKey and return correct instance or error
func NewRSATokenParser(publicKey []byte) (*RSATokenParser, error) {
	k, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}
	return &RSATokenParser{
		publicKey: k,
	}, nil
}

// RSATokenParser for parsing JWT token
// Implements TokenParser
type RSATokenParser struct {
	publicKey *rsa.PublicKey
}

// ParseCtxWithClaims takes context, parse JWT token, convert context and, if token valid, calls f with converted context and JWT Claims
func (p *RSATokenParser) ParseCtxWithClaims(context Context, f CtxClaimsFunc) error {
	c, err := p.Parse(context.Token)
	if err != nil {
		return err
	}
	ctx, cancel := context.ToContext()
	defer cancel()
	return f(ctx, *c)
}

// ParseWithClaims takes t, parse JWT token and, if token valid, calls f with JWT Claims
func (p *RSATokenParser) ParseWithClaims(t Token, f ClaimsFunc) error {
	c, err := p.Parse(t)
	if err != nil {
		return err
	}
	return f(*c)
}

// Parse JWT token and return Claims or error
func (p *RSATokenParser) Parse(t Token) (*Claims, error) {
	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(string(t), &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}
		return p.publicKey, nil
	})
	if token != nil && token.Valid {
		userID, err := strconv.ParseInt(claims.Subject, 10, 0)
		if err != nil {
			return nil, ErrInvalidToken
		}
		return &Claims{
			UserID: UserID(userID),
		}, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrInvalidToken
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrTokenExpired
		} else {
			return nil, ve
		}
	} else {
		return nil, err
	}
}

// TokenSigner interface for JWT token signing and point for mocking (see ctxtgtest subpackage)
type TokenSigner interface {
	Sign(c Claims, timeout time.Duration) (Token, error)
}

// NewRSATokenSigner parse privateKey and return correct instance
func NewRSATokenSigner(privateKey []byte) (*RSATokenSigner, error) {
	k, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}
	return &RSATokenSigner{
		privateKey: k,
	}, nil

}

// RSATokenSigner for signing JWT tokens
// Implements TokenSigner
type RSATokenSigner struct {
	privateKey *rsa.PrivateKey
}

// Sign and encode c with timeout, returns signed Token or error
func (r *RSATokenSigner) Sign(c Claims, timeout time.Duration) (Token, error) {
	t, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   strconv.FormatInt(int64(c.UserID), 10),
		ExpiresAt: timeNowFunc().Add(timeout).Unix(),
	}).SignedString(r.privateKey)
	return Token(t), err
}
