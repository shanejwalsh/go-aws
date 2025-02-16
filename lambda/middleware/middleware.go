package middleware

import (
	"fmt"
	"lambda-func/types"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const AUTH_HEADER = "Authorization"

func ValidateJWTMiddleware(next types.Next) types.Next {

	return func(req types.Req) (types.Res, error) {
		token := extractTokenFromHeaders(req.Headers)

		if token == "" {
			return types.Res{
				Body:       "No Auth Token",
				StatusCode: http.StatusUnauthorized,
			}, fmt.Errorf("No auth header found")
		}

		claims, err := parseToken(token)

		if err != nil {
			if token == "" {
				return types.Res{
					Body:       "Not Authorized",
					StatusCode: http.StatusUnauthorized,
				}, fmt.Errorf("not authorized")
			}
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return types.Res{
				Body:       "Token Expired",
				StatusCode: http.StatusUnauthorized,
			}, fmt.Errorf("expired token")
		}

		return next(req)

	}

}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader := headers[AUTH_HEADER]
	elements := strings.Split(authHeader, "Bearer ")

	if len(elements) != 2 {
		return ""
	}

	return elements[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(jwt *jwt.Token) (interface{}, error) {
		secret := "MY_SECRET"
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("token not ok %w", err)
	}

	return claims, nil

}
