# Setup Swagger

### Install Swag CLI
```
go install github.com/swaggo/swag/cmd/swag@latest
```

### Install Swag libraries
```
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
```

### Checking Swag version
```
swag --version
```

### Add headers to the main.go
```
package main

// @title Todo API
// @version 1.0
// @description Simple Todo REST API
// @BasePath /
```

### Generate swagger folder
```
swag init -g cmd/api/main.go
```
