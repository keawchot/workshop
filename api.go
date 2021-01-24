package workshop

import (
	"workshop/auth"
	"workshop/db"
	_ "workshop/docs"
	"workshop/user"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// StartServer : Create new Router with Gin
func StartServer() {
	router := echo.New()
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	router.Validator = &Validator{validator: validator.New()}
	// ===== Middlewares

	// ===== Prefix of all routes
	publicRoute := router.Group("/api/v1")
	// ===== Initial resource from MongoDB
	resource, err := db.CreateResource()

	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	// ===== Add routes of users
	user.NewUserAPI(publicRoute, resource)
	auth.NewAuthAPI(publicRoute, resource)

	// ===== Start server
	router.Start(":8080") // listen and serve on 0.0.0.0:8080
}
