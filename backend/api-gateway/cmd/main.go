package main

import "api-gateway/internal/app"

// @title           API Gateway
// @version         1.0
// @description     Central entry point for microservice design
// @description     Proxy requests to the auth-service and user-service
// @description     Enables rate limiting, CORS middleware, Auth middleware

// @contact.name    Contact Contact
// @contact.email   test@example.com

// @host            localhost
// @BasePath        /api

// @tag.name auth
// @tag.description Authentication and authorization (proxied to auth-service)

// @tag.name user
// @tag.description User management (proxied in the user-service)

// @tag.name health
// @tag.description Health checks
func main() {
	app.Run()
}
