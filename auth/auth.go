package auth

import (
	"net/http"
	"time"
	"workshop/db"
	j "workshop/jwt"
	u "workshop/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// NewUserAPI to create the router of user
func NewAuthAPI(app *echo.Group, resource *db.Resource) {
	// Create repository
	repository := u.NewUserRepository(resource)
	app.POST("/login", handleLogin(repository))
}

func handleLogin(repository u.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		user, err := repository.GetByEmailAndPassword(email, password)

		if email == "" || password == "" {
			return echo.ErrBadRequest
		}

		// Throws unauthorized error
		if err != nil || user == nil {
			return echo.ErrUnauthorized
		}

		// Set custom claims
		claims := &j.JwtCustomClaims{
			user.Id.String(),
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	}
}
