package handlers

import (
	"context"
	"cs-go/db"
	"cs-go/models"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type contextKey string

const userContextKey = contextKey("userID")

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Obtener el encabezado Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendError(rw, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Verificar el formato del token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			sendError(rw, http.StatusUnauthorized, "Invalid token format")
			return
		}

		// Parsear el token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			sendError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Verificar si el token es válido
		if !token.Valid {
			sendError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extraer el email del usuario del token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}

		userEmail, ok := claims["sub"].(string)
		if !ok {
			sendError(rw, http.StatusUnauthorized, "Invalid token: user email not found")
			return
		}

		// Obtener el ID de usuario basado en el email
		userID, err := getUserIDByEmail(userEmail)
		if err != nil {
			sendError(rw, http.StatusInternalServerError, "Error retrieving user ID")
			return
		}

		// Agregar el ID de usuario al contexto
		ctx := context.WithValue(r.Context(), userContextKey, userID)
		r = r.WithContext(ctx)

		// Si todo está bien, continuar con la siguiente función
		next.ServeHTTP(rw, r)
	})
}

func getUserIDByEmail(email string) (uint, error) {
	var user models.User
	// Buscar el usuario por email en la base de datos
	result := db.Database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, result.Error
	}
	return user.ID, nil
}
