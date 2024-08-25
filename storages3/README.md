# Module `Storage`

### Storage APIs
```bash
GET /api/v1/storage/presigned-url
PUT /api/v1/storage/legitimize-files	
``` 

### Usage
Install
```bash
go get -u github.com/gflydev/modules/storages3@v1.0.0
```

File `main.go`
```go
import (
    _ "github.com/gflydev/storage/s3"
)
```

File `api_routes.go`
```go
import "github.com/gflydev/modules/storages3/api"

// `API` Router
g.Group(prefixAPI, func(apiRouter *core.Group) {
    /* ============================ Storage Group ==========================================*/
    apiRouter.Group("/storage", func(uploadGroup *core.Group) {
        uploadGroup.GET("/presigned-url", api.NewPresignedURLApi())      // Get presigned URL
        uploadGroup.PUT("/legitimize-files", api.NewLegitimizeFileApi()) // Legitimize uploaded file
    })
})
```
