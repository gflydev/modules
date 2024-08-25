# Module `Storage`

### Storage APIs
```bash
GET /api/v1/storage/presigned-url
POST /api/v1/storage/uploads
PUT /api/v1/storage/uploads/{file_name}
PUT /api/v1/storage/legitimize-files
``` 

### Usage
Install
```bash
go get -u github.com/gflydev/modules/storage@v1.0.1
```

File `main.go`
```go
import (
    _ "github.com/gflydev/cache/redis"
    _ "github.com/gflydev/storage/local"
)
```

File `api_routes.go`
```go
import "github.com/gflydev/modules/storage/api"

// `API` Router
g.Group(prefixAPI, func(apiRouter *core.Group) {
    /* ============================ Storage Group ==========================================*/
    apiRouter.Group("/storage", func(uploadGroup *core.Group) {
        uploadGroup.GET("/presigned-url", api.NewPresignedURLApi())      // Get presigned URL
        uploadGroup.POST("/uploads", api.NewUploadApi())                 // Upload files to server by Field Form
        uploadGroup.PUT("/uploads/{file_name}", api.NewUploadFileApi())  // Upload a file to server via Body (Binary)
        uploadGroup.PUT("/legitimize-files", api.NewLegitimizeFileApi()) // Legitimize uploaded file
    })
})
```
