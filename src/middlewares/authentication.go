package middlewares

import "net/http"


type AuthMiddleware struct {
	apiKey string
}

func NewAuthMiddleware(apiKey string) *AuthMiddleware {
	return &AuthMiddleware{apiKey: apiKey}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("apikey") != m.apiKey {
			http.Error(w, "não autorizado", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}