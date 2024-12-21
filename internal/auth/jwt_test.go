package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateJWT(t *testing.T) {
	const (
		secret   = "test-secret"
		issuer   = "tainha"
		audience = "tainha-api"
	)

	tests := []struct {
		name           string
		headers        map[string]string
		claims         map[string]interface{}
		want           int
		expectedHeader map[string]string
	}{
		{
			name:    "no secret header",
			headers: map[string]string{},
			want:    http.StatusUnauthorized,
		},
		{
			name: "valid token with user claims",
			headers: map[string]string{
				HeaderJWTSecret:     secret,
				HeaderAuthorization: "Bearer ",
			},
			claims: map[string]interface{}{
				"username": "john.doe",
				"email":    "john@example.com",
				"role":     "admin",
			},
			want: http.StatusOK,
			expectedHeader: map[string]string{
				"X-Username": "john.doe",
				"X-Email":    "john@example.com",
				"X-Role":     "admin",
			},
		},
		{
			name: "valid token with standard claims",
			headers: map[string]string{
				HeaderJWTSecret:     secret,
				HeaderJWTIssuer:     issuer,
				HeaderJWTAudience:   audience,
				HeaderAuthorization: "Bearer ",
			},
			claims: map[string]interface{}{
				"iss":      issuer,
				"aud":      audience,
				"sub":      "user123",
				"username": "john.doe",
			},
			want: http.StatusOK,
			expectedHeader: map[string]string{
				"X-Sub":      "user123",
				"X-Username": "john.doe",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token if claims are provided
			if tt.claims != nil {
				token := generateToken(secret, tt.claims)
				tt.headers[HeaderAuthorization] += token
			}

			handler := ValidateJWT(secret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check that JWT configuration headers were removed
				if r.Header.Get(HeaderJWTSecret) != "" {
					t.Error("JWT secret header was not removed")
				}

				// Check that claims were properly forwarded as headers
				for key, value := range tt.expectedHeader {
					if got := r.Header.Get(key); got != value {
						t.Errorf("Expected header %s=%s, got %s", key, value, got)
					}
				}

				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.want {
				t.Errorf("ValidateJWT() status = %v, want %v", rr.Code, tt.want)
			}
		})
	}
}

func generateToken(secret string, claims map[string]interface{}) string {
	if claims == nil {
		claims = map[string]interface{}{
			"exp": time.Now().Add(time.Hour).Unix(),
		}
	}
	// Ensure expiration is set
	if _, ok := claims["exp"]; !ok {
		claims["exp"] = time.Now().Add(time.Hour).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
