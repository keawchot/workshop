package user

import (
	"net/http"
	"workshop/db"
	j "workshop/jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewUserAPI to create the router of user
func NewUserAPI(app *echo.Group, resource *db.Resource) {
	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &j.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}

	// Create repository
	repository := NewUserRepository(resource)
	app.GET("/users", handleGetUsers(repository), middleware.JWTWithConfig(config))
	app.GET("/users/me", handleGetUserByID(repository), middleware.JWTWithConfig(config))
	app.PUT("/users/me", handleUpdateUser(repository), middleware.JWTWithConfig(config))
	app.POST("/users", handleCreateNewUser(repository))
}

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age" validate:"required"`
}

type UserUpdateRequest struct {
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age" validate:"required"`
}

// Handlers

// GetUsers godoc
// @Summary Retrieves users based on query
// @Description Get Users
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Success 200 {array} Users
// @Router /api/v1/users [get]
func handleGetUsers(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		users, err := repository.GetAll()
		if err != nil {
			code = http.StatusInternalServerError
		}
		if len(users) == 0 {
			code = http.StatusNotFound
		}
		return c.JSON(code, users)
	}
}

// @Summary Get user
// @Description Get User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer"
// @Success 200 {object} User
// @Router /api/v1/users/me [get]
func handleGetUserByID(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*j.JwtCustomClaims)
		id := claims.UserID
		user, err := repository.GetByID(id)
		response := map[string]interface{}{
			"user": user,
			"err":  getErrorMessage(err),
		}
		return c.JSON(code, response)
	}
}

// @Summary Create User
// @Description Create User
// @Accept  json
// @Produce  json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string true "Name"
// @Param age formData number true "Age"
// @Success 200 {object} User
// @Router /api/v1/users [post]
func handleCreateNewUser(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		newUser := UserRequest{}
		if err := c.Bind(&newUser); err != nil {
			return c.JSON(http.StatusBadRequest, getErrorMessage(err))
		}
		if err := c.Validate(&newUser); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, getErrorMessage(err))
		}
		// Validate input !!!

		existUser, err := repository.GetByEmail(newUser.Email)

		if existUser != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Exising email!")
		}

		// Create data in database
		user, err := repository.CreateOne(newUser)
		response := map[string]interface{}{
			"user": user,
			"err":  getErrorMessage(err),
		}
		return c.JSON(code, response)
	}
}

// @Summary Update User
// @Description Update User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer"
// @Param password formData string true "Password"
// @Param name formData string true "Name"
// @Param age formData number true "Age"
// @Success 200 {object} User
// @Router /api/v1/users/me [put]
func handleUpdateUser(repository Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		code := http.StatusOK
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*j.JwtCustomClaims)
		id := claims.UserID
		newUser := UserUpdateRequest{}
		if err := c.Bind(&newUser); err != nil {
			return c.JSON(http.StatusBadRequest, getErrorMessage(err))
		}
		if err := c.Validate(&newUser); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, getErrorMessage(err))
		}
		// Validate input !!!

		// Create data in database
		user, err := repository.UpdateUser(id, newUser)
		response := map[string]interface{}{
			"user": user,
			"err":  getErrorMessage(err),
		}
		return c.JSON(code, response)
	}
}

func getErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
