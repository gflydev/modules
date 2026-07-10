package services

import (
	"errors"
	"fmt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/auth"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"strings"
	"time"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      int
	Credentials core.Data
	Expires     int64
}

// GenerateTokens func for generate a new Access & Refresh tokens.
func GenerateTokens(id string, credentials []string) (*auth.Token, error) {
	// Generate JWT Access token.
	accessToken, err := generateAccessToken(id, credentials)
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	// Generate JWT Refresh token.
	refreshToken, err := generateRefreshToken()
	if err != nil {
		// Return refresh generation error.
		return nil, err
	}

	return &auth.Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateAccessToken(id string, credentials []string) (string, error) {
	// Get secret key from .env file.
	secret := utils.Getenv(auth.SecretKey, "")

	// Set expired minutes count for a secret key from .env file.
	ttlMinutes := utils.Getenv(auth.TtlMinutes, 0)

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["id"] = id
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(ttlMinutes)).Unix()

	// Set private token credentials:
	for _, credential := range credentials {
		claims[credential] = true
	}

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

// IsValidRefreshToken func for parse second argument from refresh token.
// A refresh token is valid is not expired.
func IsValidRefreshToken(refreshToken string) bool {
	tokenString := strings.Split(refreshToken, ".")
	if len(tokenString) < 2 {
		return false
	}
	expires, err := strconv.ParseInt(tokenString[1], 10, 64)
	if err != nil {
		log.Infof("parse refresh token error %v", err)

		return false
	}

	if expires < time.Now().Unix() {
		log.Info("refresh token expired")

		return false
	}

	return true
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
		Credentials: make(core.Data),
		Expires:     int64(expiresFloat),
	}, nil
}

func generateRefreshToken() (string, error) {
	hash := utils.Sha256(utils.Getenv(auth.RefreshKey, "") + time.Now().String())

	// Get expired days for refresh key from .env file.
	overDays := utils.Getenv(auth.TtlOverDays, 0)

	// Create expiration time.
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(overDays*24)).Unix())

	// Create a new refresh token (sha256 string with salt + expire time).
	t := hash + "." + expireTime

	return t, nil
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
	return []byte(os.Getenv(auth.SecretKey)), nil
}
