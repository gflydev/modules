# Module `Storage`

### Auth APIs
```bash
POST /api/v1/auth/signin
POST /api/v1/auth/signup
DELETE /api/v1/auth/signout
PUT /api/v1/auth/refresh
``` 

### Usage
Install
```bash
go get -u github.com/gflydev/modules/jwt@v1.0.0
```

File `main.go`
```go
import (
    _ "github.com/gflydev/cache/redis"
    _ "github.com/gflydev/db/psql"
    mb "github.com/gflydev/db"	
    "github.com/gflydev/modules/jwt"
    "github.com/gflydev/modules/jwt/api"
)
```

File `api_routes.go`
```go
// `API` Router
g.Group(prefixAPI, func(apiRouter *core.Group) {
    apiRouter.Use(jwt.New(
        prefixAPI+"/auth/signin",
        prefixAPI+"/auth/signup",
        prefixAPI+"/auth/refresh",
    ))

    /* ============================ Auth Group ===================================*/
    apiRouter.Group("/auth", func(authGroup *core.Group) {
        authGroup.POST("/signin", api.NewSignInApi())
        authGroup.DELETE("/signout", api.NewSignOutApi())
        authGroup.POST("/signup", api.NewSignUpApi())
        authGroup.PUT("/refresh", api.NewRefreshTokenApi())
    })
})
```

### Tables

On `PostgreSQL`
```sql
-- -----------------------------------------------------
-- Table `users`
-- -----------------------------------------------------
CREATE TYPE user_status AS ENUM ('pending', 'active', 'blocked');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR (255) NOT NULL UNIQUE,
    password VARCHAR (255) NOT NULL,
    fullname VARCHAR (255) NULL,
    phone VARCHAR(20) NULL,
    token VARCHAR (100) NULL,
    status user_status DEFAULT 'pending',
    avatar VARCHAR (255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    verified_at TIMESTAMP NULL,
    blocked_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    last_access_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id);
CREATE UNIQUE INDEX email_users ON users (email ASC);
```

On `MySQL`
```sql
-- -----------------------------------------------------
-- Table users
-- -----------------------------------------------------
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR (255) NOT NULL UNIQUE,
    password VARCHAR (255) NOT NULL,
    fullname VARCHAR (255) NULL,
    phone VARCHAR(20) NULL,
    token VARCHAR (100) NULL,
    status ENUM('pending', 'active', 'blocked') NOT NULL DEFAULT 'pending',
    avatar VARCHAR (255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    verified_at TIMESTAMP NULL,
    blocked_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    last_access_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id);
CREATE UNIQUE INDEX email_users ON users (email ASC);
```