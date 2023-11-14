package middleware

import (
	"time"

	"github.com/dedihartono801/chat-realtime/pkg/config"
	"github.com/dedihartono801/chat-realtime/pkg/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Define the claims structure for the JWT
type Claims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

// Define a function for generating a new JWT
func GenerateToken(id int64) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Define a middleware for verifying JWT authentication and expiration
func AuthUser(c *fiber.Ctx) error {
	// Get the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helpers.CustomResponse(c, nil, "Unauthorized", 401)
	}

	// Verify that the Authorization header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return helpers.CustomResponse(c, nil, "Invalid format authorization", 401)
	}

	// Parse the JWT from the Authorization header
	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return helpers.CustomResponse(c, nil, "Invalid signature", 401)
		}
		return helpers.CustomResponse(c, nil, "Unauthorized", 401)
	}

	// Check if the JWT has expired
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return helpers.CustomResponse(c, nil, "Unauthorized", 401)
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return helpers.CustomResponse(c, nil, "Expired token", 401)
	}

	// Set the user ID in the context for future requests
	c.Locals("ID", claims.ID)

	// Call the next middleware in the chain
	return c.Next()
}
