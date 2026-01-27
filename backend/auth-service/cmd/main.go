package main

import "auth-service/internal/app"

// @title           Auth Service API
// @version         1.0
// @description     API for authentication JWT tokens.
// @contact.name    Contact Contact
// @contact.email   contact@example.com

// @host            localhost:8081
// @BasePath        /api/auth

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {Your JWT token}" to authenticate.
func main() {
	app.Run()
}
