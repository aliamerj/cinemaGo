package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Middleware for authentication
func AuthenticateMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		tokenHeader := r.Header.Get("Authorization")

		// Split token to remove Bearer prefix
		splitToken := strings.Split(tokenHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid/Malformed token", http.StatusForbidden)
			return
		}

		tokenPart := splitToken[1] // Grab the token part, what we are truly interested in

		// Validate token
		token, err := jwt.ParseWithClaims(tokenPart, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("mySecretKey"), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)

		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// Now add the user ID to the request context
		r = r.WithContext(context.WithValue(r.Context(), "userID", claims.Subject))

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
