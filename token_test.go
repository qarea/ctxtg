package ctxtg

import (
	"context"
	"crypto/rsa"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestParseCtxWithClaims(t *testing.T) {
	c := jwt.StandardClaims{
		Subject:   "3",
		ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)
	p := testRSATokenParser(t)

	contexttg := Context{
		Token: str,
	}

	err := p.ParseCtxWithClaims(contexttg, func(ctx context.Context, claims Claims) error {
		if !reflect.DeepEqual(FromContext(ctx), contexttg) {
			t.Error("Invalid context passed")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestParseCtxWithClaimsErr(t *testing.T) {
	c := jwt.StandardClaims{
		Subject:   "3",
		ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)
	p := testRSATokenParser(t)

	contexttg := Context{
		Token: str,
	}

	testErr := errors.New("claims err")

	err := p.ParseCtxWithClaims(contexttg, func(ctx context.Context, claims Claims) error {
		if !reflect.DeepEqual(FromContext(ctx), contexttg) {
			t.Error("Invalid context passed")
		}
		return testErr
	})
	if err != testErr {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestParseCtxWithClaimsParserErr(t *testing.T) {
	token := Token("dsda")
	p := testRSATokenParser(t)

	contexttg := Context{
		Token: token,
	}

	err := p.ParseCtxWithClaims(contexttg, func(ctx context.Context, claims Claims) error {
		t.Error("Should not be called")
		return nil
	})
	if err != ErrInvalidToken {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestParseWithClaims(t *testing.T) {
	c := jwt.StandardClaims{
		Subject:   "3",
		ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)
	p := testRSATokenParser(t)

	err := p.ParseWithClaims(str, func(claims Claims) error {
		if c.Subject != strconv.FormatInt(int64(claims.UserID), 10) {
			t.Errorf("invalid claims %v", claims)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestParseWithClaimsParserErr(t *testing.T) {
	token := Token("dasfas")
	p := testRSATokenParser(t)

	err := p.ParseWithClaims(token, func(claims Claims) error {
		t.Error("Should not be called")
		return nil
	})
	if err != ErrInvalidToken {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestParseWithClaimsErr(t *testing.T) {
	c := jwt.StandardClaims{
		Subject:   "3",
		ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)
	p := testRSATokenParser(t)
	testErr := errors.New("ClaimsFuncErr")

	err := p.ParseWithClaims(str, func(claims Claims) error {
		if c.Subject != strconv.FormatInt(int64(claims.UserID), 10) {
			t.Errorf("invalid claims %v", claims)
		}
		return testErr
	})
	if err != testErr {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestRSATokenParserTest(t *testing.T) {
	c := jwt.StandardClaims{
		Subject:   "3",
		ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)
	p := testRSATokenParser(t)

	claims, err := p.Parse(str)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if c.Subject != strconv.FormatInt(int64(claims.UserID), 10) {
		t.Errorf("invalid claims %v", claims)
	}
}

func TestRSATokenParserInvalidKey(t *testing.T) {
	_, err := NewRSATokenParser([]byte("invalidkey"))
	if err == nil {
		t.Errorf("Should return err")
	}
}

func TestRSATokenParserExpiredTest(t *testing.T) {
	c := jwt.StandardClaims{
		Issuer:    "3",
		ExpiresAt: time.Now().Add(-5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)

	p := testRSATokenParser(t)
	claims, err := p.Parse(str)
	if claims != nil {
		t.Errorf("Claims should be empty %v", claims)
	}
	if err == nil {
		t.Error("Expected error")
	}
	if err != ErrTokenExpired {
		t.Errorf("TokenExpired error expected %v %T", err, err)
	}
}

func TestRSATokenParserInvalidMethod(t *testing.T) {
	c := jwt.StandardClaims{
		Issuer:    "3",
		ExpiresAt: time.Now().Add(-5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, c)
	str, err := token.SigningString()
	if err != nil {
		t.Fatal(err)
	}

	p := testRSATokenParser(t)

	claims, err := p.Parse(Token(str))
	if claims != nil {
		t.Errorf("Claims should be empty %v", claims)
	}
	if err == nil {
		t.Error("Expected error")
	}
	if err != ErrInvalidToken {
		t.Errorf("Invalid token error expected %v %T", err, err)
	}
}

func TestRSATokenParserMalformed(t *testing.T) {
	p := testRSATokenParser(t)
	claims, err := p.Parse("invalidkey.dsads.dsad")
	if claims != nil {
		t.Errorf("Claims should be empty %v", claims)
	}
	if err == nil {
		t.Error("Expected error")
	}
	if err != ErrInvalidToken {
		t.Errorf("Invalid token error expected %v %T", err, err)
	}
}

func TestRSATokenStringUserID(t *testing.T) {
	c := jwt.MapClaims{
		"iss": "id3",
		"exp": time.Now().Add(5 * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	str := signToken(t, token)

	p := testRSATokenParser(t)
	claims, err := p.Parse(str)
	if claims != nil {
		t.Errorf("UserID should be empty %v", claims)
	}
	if err == nil {
		t.Error("Expected error")
	}
	if err != ErrInvalidToken {
		t.Errorf("TokenExpired error expected %v %T", err, err)
	}
}

func TestRSATokenSigner(t *testing.T) {
	now, f := testTime()
	defer f()
	tokenDuration := 5 * time.Second
	s := testRSATokenSigner(t)
	claims := Claims{
		UserID: 1,
	}
	token, err := s.Sign(claims, tokenDuration)
	if err != nil {
		t.Error(t)
	}
	c, expiredAt := parseToken(t, token)
	if claims.UserID != c.UserID {
		t.Error("User id invalid", c.UserID)
	}
	if now.Add(tokenDuration).Equal(time.Unix(expiredAt, 0)) {
		t.Error("Invalid expiredAt expired", expiredAt)
	}
}

func TestRSATokenSignerInvalidKey(t *testing.T) {
	_, err := NewRSATokenSigner([]byte("invalidkey"))
	if err == nil {
		t.Errorf("Should be error")
	}
}

func testTime() (time.Time, func()) {
	testTime := time.Now().Add(10 * time.Second)
	timeNowFunc = func() time.Time {
		return testTime
	}
	return testTime, func() {
		timeNowFunc = time.Now
	}
}

func parseToken(t *testing.T, token Token) (Claims, int64) {
	var c jwt.StandardClaims
	_, err := jwt.ParseWithClaims(string(token), &c, func(*jwt.Token) (interface{}, error) {
		return testPublicKey(t), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	userID, err := strconv.ParseInt(c.Subject, 10, 0)
	if err != nil {
		t.Fatal(t)
	}
	return Claims{
		UserID: UserID(userID),
	}, c.ExpiresAt
}

func signToken(t *testing.T, jt *jwt.Token) Token {
	token, err := jt.SignedString(testPrivateKey(t))
	if err != nil {
		t.Fatal(err)
	}
	return Token(token)
}

func testPrivateKey(t *testing.T) *rsa.PrivateKey {
	k, err := jwt.ParseRSAPrivateKeyFromPEM(privateRSA)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func testPublicKey(t *testing.T) *rsa.PublicKey {
	k, err := jwt.ParseRSAPublicKeyFromPEM(publicRSA)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func testRSATokenSigner(t *testing.T) *RSATokenSigner {
	s, err := NewRSATokenSigner(privateRSA)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func testRSATokenParser(t *testing.T) *RSATokenParser {
	p, err := NewRSATokenParser(publicRSA)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

var privateRSA = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOQIBAAJAcr5bdI/2NZ2DpMwh2J945xAPGkBkrCGmSuAy9SqPiL46jQQvZt68
m7AxHQkG/JLhMql1xwjesoQeSoKz5LpdSwIDAQABAkA8d9aoacl9XcHnUfAwQXIs
ioj6855aG+2PnfEcxE4Z6DYh68JJiFgXVqXoLeL7DTKMffXFWGpZfPkwF/oNQzsJ
AiEAzj4hfWNLNTj+QS6wlFXeUWpD5OvhX2PN0AA0tSN4L9cCIQCObRlaOFUoI6mW
0FHT09Kqh2Eb7JROKC4uOzvx7CEfrQIgW/7Iz3ZoCLCIcSjTaQc4aJZ+/HDfEb6i
AmLlH9tXc/cCIAW93DHI55XwqhuMVmAlv+5j+sQ3a1sjP4lZlfcQv90ZAiEAgrk1
SO1DU8Q5HedenIIZEp8BF2yXCkuWdRKTbSCt2ZY=
-----END RSA PRIVATE KEY-----`)

var publicRSA = []byte(`-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAcr5bdI/2NZ2DpMwh2J945xAPGkBkrCGm
SuAy9SqPiL46jQQvZt68m7AxHQkG/JLhMql1xwjesoQeSoKz5LpdSwIDAQAB
-----END PUBLIC KEY-----`)
