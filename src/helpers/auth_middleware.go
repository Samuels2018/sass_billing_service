package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	useranme string `json:"username"`
	jwt.RegisteredClaims
}

var db *sql.DB

func validateToken(tokenString string, jwtSecret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if Claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return Claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func AuthMiddleware(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	jwtSecret := []byte(os.Getenv(os.Getenv("JWT_SECRET")))
	if jwtSecret == nil {
		log.Fatal("JWT_SECRET not set in .env file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	Claims, err := validateToken(tokenString, jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	// Almacenar los claims en el contexto local de Fiber
	c.Locals("user", Claims.useranme)
	log.Println("User authenticated:", Claims.useranme)

	// Continuar con el siguiente middleware/handler
	return c.Next()

}
