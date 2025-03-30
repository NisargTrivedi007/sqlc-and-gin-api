package jwt_code

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Default environment variable name for JWT key
const DefaultJWTKeyEnv = "JWT_Key"

// JWT struct to hold the secret key

// LoadJWTKey loads the JWT key from a .env file and returns it
func LoadJWTKey(envFile string, keyName string) (string, error) {
	if keyName == "" {
		keyName = DefaultJWTKeyEnv
	}

	// Load env file if provided
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			return "", fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Get the key from environment
	key := os.Getenv(keyName)
	if key == "" {
		return "", fmt.Errorf("JWT key not found in environment variable %s", keyName)
	}

	return key, nil
}

func GenerateToken(userID int, username string) (string, error) {
	tokenSecret, err := LoadJWTKey(".env", DefaultJWTKeyEnv)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	tokenSecret, err := LoadJWTKey(".env", DefaultJWTKeyEnv)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Verify that the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// GetUserIDFromToken extracts the user ID from a validated token
func GetUserIDFromToken(token *jwt.Token) (int, error) {
	// tokenSecret := os.Getenv(DefaultJWTKeyEnv)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("could not parse claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id claim not found or invalid format")
	}

	return int(userID), nil
}
