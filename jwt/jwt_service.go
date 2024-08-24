package jwt

import (
	"fmt"
	"github.com/gflydev/cache"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/db/null"
	"github.com/gflydev/modules/jwt/dto"
	"github.com/gflydev/modules/jwt/model"
	"github.com/gflydev/modules/jwt/repository"
	_utils "github.com/gflydev/modules/jwt/utils"
	"strconv"
	"time"
)

// SignIn login app.
func SignIn(signIn *dto.SignIn) (*Tokens, error) {
	// Get user by email.
	user := repository.Pool.GetUserByEmail(signIn.Username)
	if user == nil {
		return nil, errors.New("Invalid email address or password")
	}
	// Compare given user password with stored in found user.
	isValidPassword := _utils.ComparePasswords(user.Password, signIn.Password)
	if !isValidPassword {
		return nil, errors.New("Invalid email address or password")
	}

	userIDStr := strconv.Itoa(user.ID)
	// Generate a new pair of access and refresh tokens.
	tokens, err := GenerateTokens(userIDStr, make([]string, 0))
	if err != nil {
		log.Errorf("Error while generating tokens %q", err)
		return nil, err
	}

	// Set expired days from .env file
	ttlDays := utils.Getenv("JWT_TTL_OVER_DAYS", 0) // 7 days by default

	// Save refresh token to Redis.
	expiredTime := time.Duration(ttlDays*24*3600) * time.Second // 604 800 seconds = 7 days
	if err = cache.Set(userIDStr, tokens.Refresh, expiredTime); err != nil {
		log.Errorf("Error while caching to token to Redis %q", err)
		return nil, err
	}

	return tokens, nil
}

// SignUp register new user.
func SignUp(signUp *dto.SignUp) (*model.User, error) {
	userEmail := repository.Pool.GetUserByEmail(signUp.Email)
	if userEmail != nil {
		return nil, errors.New("user with the given email address already exists")
	}

	// Create a new user struct.
	user := &model.User{}

	// If user pass-in `status`, replace status variable
	status := model.UserStatusActive
	if signUp.Status != "" {
		status = signUp.Status
	}

	// Set initialized default data for user
	user.Email = signUp.Email
	user.Password = _utils.GeneratePassword(signUp.Password)
	user.Fullname = signUp.Fullname
	user.Phone = signUp.Phone
	user.Token = null.String("")
	user.Status = status
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastAccessAt = null.NowTime()

	// Create a new user with validated data.
	err := repository.Pool.CreateUser(user)
	if err != nil {
		log.Errorf("Error while creating new user %q with data '%v'", err, user)
		return nil, errors.New("error occurs while signup user")
	}

	return user, nil
}

// SignOut function takes in jwtToken string, utils.ExtractTokenMetadata extract access token metadata
// to get a userID which is the key that store refresh token in the Redis Caching
// then delete refresh token from the Redis
// and DeleteJWTToken will delete access token
// by send it to black-list (middleware will handle invalid token in blacklist).
func SignOut(jwtToken string) error {
	// Extract access token metadata
	claims, err := ExtractTokenMetadata(jwtToken)
	if err != nil {
		log.Errorf("Error while logging out %q", err)
		return errors.New("logout error")
	}

	userID := strconv.Itoa(claims.UserID)

	// Delete refresh token from Redis.
	if err = cache.Del(userID); err != nil {
		log.Errorf("Error while delete refresh token from Redis %q", err)
		return err
	}

	// Delete access token by send it to black-list
	DeleteJWTToken(jwtToken)
	return nil
}

// RefreshJWTToken function to refresh JWT token from user.
func RefreshJWTToken(jwtToken, refreshToken string) (*Tokens, error) {
	// Get claims from JWT.
	claims, err := ExtractTokenMetadata(jwtToken)
	if err != nil {
		log.Errorf("Error while extracting token metadata %q", err)
		return nil, errors.New("refresh token error")
	}
	// Define user ID.
	userID := claims.UserID
	userIDStr := strconv.Itoa(userID)

	// Get refresh token from Redis.
	val, err := cache.Get(userIDStr)
	if err != nil {
		log.Errorf("Error while getting refresh token from Redis %q", err)
		return nil, errors.New("refresh token error")
	}

	if refreshToken != val {
		log.Errorf("Mismatch refresh token '%s' vs input refresh token '%s'", refreshToken, val)

		return nil, errors.New("refresh token mismatch")
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := GenerateTokens(userIDStr, make([]string, 0))
	if err != nil {
		log.Errorf("Error while generating JWT Tokens")
		return nil, errors.New("refresh token error")
	}

	// Set expired days from .env file.
	ttlDays := utils.Getenv("JWT_TTL_OVER_DAYS", 0)
	duration := time.Duration(ttlDays*7*24*3600) * time.Second

	// Update refresh token to Redis.
	if err = cache.Set(userIDStr, tokens.Refresh, duration); err != nil {
		log.Errorf("Refresh token error '%v'", err)

		return nil, errors.New("refresh token error")
	}

	// Delete JWT token by sending it to blacklist
	DeleteJWTToken(jwtToken)

	return tokens, nil
}

// GetUserByToken returns User by JWT token
func GetUserByToken(jwtToken string) *model.User {
	// Get claims from JWT.
	claims, err := ExtractTokenMetadata(jwtToken)
	if err != nil {
		log.Errorf("Get user from JWT token error '%v'", err)

		return nil
	}

	// Define user ID.
	userID := claims.UserID
	userIDStr := strconv.Itoa(userID)

	// Get refresh token to Redis.
	if _, err = cache.Get(userIDStr); err != nil {
		log.Errorf("Get user from JWT token error '%v'", err)

		return nil
	}

	// Get user by ID.
	return repository.Pool.GetUserByID(userID)
}

// DeleteJWTToken add jwtToken to blacklist
func DeleteJWTToken(jwtToken string) bool {
	key := fmt.Sprintf("%s:%s", utils.Getenv("JWT_BLACKLIST", ""), jwtToken)

	// Set expired minutes count for a secret key from .env file.
	ttlMinutes := utils.Getenv("JWT_TTL_MINUTES", 0)
	expiresTime := time.Duration(ttlMinutes*60) * time.Second

	// Update refresh token to Redis.
	if err := cache.Set(key, "blocked", expiresTime); err != nil {
		log.Errorf("Delete JWT token error '%v'", err)

		return false
	}

	return true
}

// IsBlockedJWTToken Check if jwtToken is locked or not
func IsBlockedJWTToken(jwtToken string) (bool, error) {
	key := fmt.Sprintf("%s:%s", utils.Getenv("JWT_BLACKLIST", ""), jwtToken)

	// Get refresh token to Redis.
	val, err := cache.Get(key)
	if err != nil {
		return false, nil
	}
	exists := val == model.UserStatusBlocked

	return exists, nil
}
