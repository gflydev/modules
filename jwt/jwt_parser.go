package jwt

import (
	"errors"
	"fmt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"strings"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      int
	Credentials core.Data
	Expires     int64
}

// ExtractToken func to get JWT from header.
func ExtractToken(c *core.Ctx) string {
	bearToken := c.Root().Request.Header.Peek(core.HeaderAuthorization)

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(utils.UnsafeStr(bearToken), " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(tokenString string) (*TokenMetadata, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Safely read the `id` claim (stored as a string).
	idStr, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("invalid token: missing id claim")
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token: bad id claim: %w", err)
	}

	// Safely read the `expires` claim (stored as a JSON number).
	expiresFloat, ok := claims["expires"].(float64)
	if !ok {
		return nil, errors.New("invalid token: missing expires claim")
	}

	return &TokenMetadata{
		UserID:      userID,
		Credentials: core.Data{},
		Expires:     int64(expiresFloat),
	}, nil
}

// verifyToken function will parse, validate and verify the signature
func verifyToken(tokenString string) (*jwt.Token, error) {
	// Pin the accepted signing algorithm to HS256 to prevent
	// algorithm-confusion / `alg: none` attacks.
	token, err := jwt.Parse(tokenString, jwtKeyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	return token, nil
}

// jwtKeyFunc will receive the parsed token and should return the cryptographic key
// for verifying the signature.
func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	// Ensure the token was signed with an HMAC method before handing back
	// the shared secret.
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
