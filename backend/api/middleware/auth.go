package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	client       *gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

type Config struct {
	KeycloakURL  string
	Realm        string
	ClientId     string
	ClientSecret string
}

func NewAuthMiddleware(cfg Config) *AuthMiddleware {
	return &AuthMiddleware{
		client:       gocloak.NewClient(cfg.KeycloakURL),
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
		realm:        cfg.Realm,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()

		// Validate access token
		_, err := m.client.RetrospectToken(ctx, token, m.clientId, m.clientSecret, m.realm)
		if err != nil {
			http.Error(w, "Token validation failed", http.StatusUnauthorized)
			return
		}

		claims, err := m.decodeClaims(token)
		if err != nil {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Check prothetic_user role
		if !m.hasProtheticUserRole(claims) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) decodeClaims(token string) (jwt.MapClaims, error) {
	parsedToken, _, err := m.client.DecodeAccessToken(context.Background(), token, m.realm)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	return claims, nil
}

func (m *AuthMiddleware) hasProtheticUserRole(claims jwt.MapClaims) bool {
	realmAccess, ok := claims["realm_access"].(map[string]interface{})
	if !ok {
		return false
	}

	roles, ok := realmAccess["roles"].([]interface{})
	if !ok {
		return false
	}

	for _, role := range roles {
		if role == "prothetic_user" {
			return true
		}
	}
	return false
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
