# Auth Module

A comprehensive authentication and authorization module for GFly framework, providing both JWT-based and session-based authentication mechanisms.

## Features

- **Multiple Authentication Methods**
  - JWT (JSON Web Token) authentication for API access
  - Session/Cookie-based authentication for web applications

- **User Management**
  - User registration (sign up)
  - User login (sign in)
  - User logout (sign out)
  - Token refresh mechanism

- **Password Management**
  - Forgot password functionality
  - Password reset with secure tokens
  - Email notifications for password changes

- **Security Features**
  - JWT token blacklisting
  - Token expiration management
  - Password hashing with bcrypt
  - Secure token generation

- **Role-Based Access Control**
  - User roles and permissions
  - Role assignment through user_roles relationship

## Installation

```bash
go get github.com/gflydev/modules/auth
```

## Dependencies

The auth module requires the following GFly packages:

- `github.com/gflydev/core` - Core framework functionality
- `github.com/gflydev/cache` - Redis cache for token management
- `github.com/gflydev/db` - Database operations
- `github.com/gflydev/http` - HTTP utilities
- `github.com/gflydev/notification` - Email notifications
- `github.com/gflydev/utils` - Utility functions
- `github.com/golang-jwt/jwt/v5` - JWT token handling

## Configuration

Add the following environment variables to your `.env` file:

```env
# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_REFRESH_KEY=your-refresh-key-here
JWT_TTL_MINUTES=60                    # Access token TTL in minutes
JWT_TTL_OVER_DAYS=7                   # Refresh token TTL in days
JWT_BLACKLIST=jwt_blacklist           # Redis key prefix for blacklist
JWT_CHECK_BLACKLIST=true              # Enable/disable blacklist checking

# API Configuration
API_PREFIX=api
API_VERSION=v1

# Auth Configuration
AUTH_LOGIN_URI=/login                 # Login page URI for session auth
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    fullname VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    token TEXT,
    status VARCHAR(20) DEFAULT 'active',
    avatar TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP,
    blocked_at TIMESTAMP,
    deleted_at TIMESTAMP,
    last_access_at TIMESTAMP
);
```

### Roles Table

```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
```

### User Roles Table

```sql
CREATE TABLE user_roles (
    user_id INTEGER REFERENCES users(id),
    role_id INTEGER REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);
```

## Usage

### 1. Register API Routes

```go
package main

import (
    "github.com/gflydev/core"
    authRoutes "github.com/gflydev/modules/auth/routes"
)

func main() {
    app := core.New()

    // Register auth routes
    apiRouter := app.Group("")
    authRoutes.RegisterApi(apiRouter)

    app.Listen(":3000")
}
```

### 2. JWT Authentication Middleware

Protect your API endpoints with JWT authentication:

```go
import (
    "github.com/gflydev/core"
    "github.com/gflydev/modules/auth/middleware"
)

// Apply JWT middleware with exclusions
apiRouter.Use(middleware.JWTAuth(
    "/api/v1/auth/signin",      // Exclude sign in
    "/api/v1/auth/signup",      // Exclude sign up
    "/api/v1/password/forgot",  // Exclude forgot password
))

// Protected routes
apiRouter.GET("/api/v1/profile", func(c *core.Ctx) error {
    // Get authenticated user from context
    user := c.GetData(http.UserKey).(models.User)
    return c.JSON(user)
})
```

### 3. Session Authentication Middleware

Protect web pages with session-based authentication:

```go
import (
    "github.com/gflydev/core"
    "github.com/gflydev/modules/auth/middleware"
)

// Apply session middleware
webRouter.Use(middleware.SessionAuth(
    "/login",           // Exclude login page
    "/register",        // Exclude register page
))

// Protected web routes
webRouter.GET("/dashboard", func(c *core.Ctx) error {
    user := c.GetData(http.UserKey).(models.User)
    return c.Render("dashboard.html", core.Data{
        "user": user,
    })
})
```

### 4. API Endpoints

The module automatically registers the following endpoints:

#### Frontend (Web) Authentication

- `POST /api/v1/frontend/auth/signin` - Web sign in (creates session)
- `DELETE /api/v1/frontend/auth/signout` - Web sign out (destroys session)

#### Backend (API) Authentication

- `POST /api/v1/auth/signin` - API sign in (returns JWT tokens)
- `POST /api/v1/auth/signup` - User registration
- `DELETE /api/v1/auth/signout` - API sign out (blacklists tokens)
- `PUT /api/v1/auth/refresh` - Refresh JWT access token

#### Password Management

- `POST /api/v1/password/forgot` - Request password reset
- `POST /api/v1/password/reset` - Reset password with token

### 5. Using Auth Services

#### Sign In

```go
import (
    "github.com/gflydev/modules/auth/dto"
    "github.com/gflydev/modules/auth/services"
)

func signIn(email, password string) (*auth.Token, error) {
    credentials := dto.SignIn{
        Username: email,
        Password: password,
    }

    tokens, err := services.SignIn(credentials)
    if err != nil {
        return nil, err
    }

    // tokens.Access contains the JWT access token
    // tokens.Refresh contains the refresh token
    return tokens, nil
}
```

#### Sign Up

```go
func signUp(email, password, fullname, phone string) (*models.User, error) {
    signUpData := dto.SignUp{
        Email:    email,
        Password: password,
        Fullname: fullname,
        Phone:    phone,
    }

    user, err := services.SignUp(signUpData)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

#### Refresh Token

```go
func refreshToken(accessToken, refreshToken string) (*auth.Token, error) {
    tokens, err := services.RefreshToken(accessToken, refreshToken)
    if err != nil {
        return nil, err
    }

    return tokens, nil
}
```

#### Forgot Password

```go
func forgotPassword(email string) error {
    forgotPW := dto.ForgotPassword{
        Username: email,
    }

    err := services.ForgotPassword(forgotPW)
    if err != nil {
        return err
    }

    // Email with reset link sent to user
    return nil
}
```

#### Reset Password

```go
func resetPassword(token, newPassword string) error {
    resetPW := dto.ResetPassword{
        Token:    token,
        Password: newPassword,
    }

    err := services.ChangePassword(resetPW)
    if err != nil {
        return err
    }

    return nil
}
```

### 6. Working with User Repository

```go
import (
    "github.com/gflydev/modules/auth/domain/repository"
)

// Get user by email
user := repository.Pool.GetUserByEmail("user@example.com")

// Get user by token
user := repository.Pool.GetUserByToken("reset_password:abc123")

// Get user by ID (using db package)
user, err := mb.GetModelByID[models.User](userID)
```

## Request/Response Examples

### Sign In Request

```bash
curl -X POST http://localhost:3000/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "access": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh": "d1a4216a226cbf75eaefc9107c2c64b6b2c0f18cd8634e3a6f495146c38e1324.1747914602"
}
```

### Sign Up Request

```bash
curl -X POST http://localhost:3000/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "SecurePass123",
    "fullname": "John Doe",
    "phone": "1234567890"
  }'
```

### Refresh Token Request

```bash
curl -X PUT http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access-token>" \
  -d '{
    "token": "<refresh-token>"
  }'
```

### Authenticated Request

```bash
curl -X GET http://localhost:3000/api/v1/protected-route \
  -H "Authorization: Bearer <access-token>"
```

## Advanced Usage

### Custom Middleware Configuration

```go
// Session authentication with manual control
webRouter.Use(middleware.SessionManipulation)

webRouter.GET("/profile",
    f.Apply(middleware.SessionAuthPage)(profileHandler))
```

### Extract User from Context

```go
func handler(c *core.Ctx) error {
    // Get authenticated user
    user := c.GetData(http.UserKey).(models.User)

    // Access user properties
    fmt.Println(user.Email, user.Fullname)

    return c.JSON(user)
}
```

### Check if Token is Blacklisted

```go
isBlocked, err := services.IsBlockedToken(jwtToken)
if err != nil || isBlocked {
    return errors.New("Token is invalid or blocked")
}
```

## User Status Types

The module defines the following user statuses:

- `active` - User account is active and can sign in
- `pending` - User account is awaiting verification
- `blocked` - User account is blocked
- `suspended` - User account is temporarily suspended

## Security Best Practices

1. **Use Strong Secret Keys**: Generate cryptographically secure random strings for `JWT_SECRET_KEY` and `JWT_REFRESH_KEY`
2. **Enable Token Blacklisting**: Set `JWT_CHECK_BLACKLIST=true` in production
3. **Configure Appropriate TTLs**: Balance security and user experience with token expiration times
4. **Use HTTPS**: Always use HTTPS in production to protect tokens in transit
5. **Validate Input**: All DTOs include validation tags - ensure validation middleware is enabled
6. **Secure Password Reset**: Reset tokens are single-use and expire after use
7. **Rate Limiting**: Consider implementing rate limiting on authentication endpoints

## Error Handling

The module returns descriptive errors for common scenarios:

- Invalid credentials: `"Invalid email address or password"`
- User not active: `"User is not activated"`
- Token expired: `"JWT token expired"`
- Token blocked: `"JWT token was blocked"`
- User not found: `"User not found"`
- Email exists: `"User with the given email address already exists"`
- Refresh token mismatch: `"Refresh token mismatch"`

## Email Notifications

The module sends email notifications for:

1. **Password Reset**: When user requests password reset
2. **Password Changed**: Confirmation email after successful password change

Configure your notification settings in the `notification` module configuration.

## Testing

Run the module tests:

```bash
cd auth
make test
```

## Project Structure

```
auth/
├── api/                    # API handlers
│   ├── signin_api.go
│   ├── signup_api.go
│   ├── signout_api.go
│   ├── refresh_token_api.go
│   ├── forgot_password_api.go
│   └── reset_password_api.go
├── domain/
│   ├── models/            # Database models
│   │   ├── user_model.go
│   │   ├── role_model.go
│   │   ├── user_roles_model.go
│   │   └── types/         # Custom types
│   └── repository/        # Data access layer
├── dto/                   # Data transfer objects
├── middleware/            # Authentication middleware
│   ├── jwt_middleware.go
│   └── session_middleware.go
├── notifications/         # Email notifications
├── request/              # Request validation
├── response/             # Response transformers
├── routes/               # Route registration
├── services/             # Business logic
│   ├── auth_services.go
│   ├── password_services.go
│   └── jwt_generator.go
├── transformers/         # Data transformers
└── init.go               # Module initialization
```

## License

This module is part of the GFly framework ecosystem.

## Support

For issues, questions, or contributions, please refer to the main GFly framework repository.
