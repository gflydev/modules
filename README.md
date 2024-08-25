# gFly modules

Common modules

## Module `JWT`

### Auth APIs
```bash
POST /api/v1/auth/signin
POST /api/v1/auth/signup
DELETE /api/v1/auth/signout
PUT /api/v1/auth/refresh
``` 

## Module `Storage`

### Local Storage APIs
```bash
POST /api/v1/storage/uploads
GET /api/v1/storage/presigned-url
PUT /api/v1/storage/uploads/{file_name}
PUT /api/v1/storage/legitimize-files
``` 

### S3 Storage APIs
```bash
GET /api/v1/storage/s3/presigned-url
PUT /api/v1/storage/s3/legitimize-files	
``` 