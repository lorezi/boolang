package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lorezi/boolang/helpers"
	"github.com/lorezi/boolang/models"
)

// Auth validate token and authorize users

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		clientToken := r.Header.Get("Authorization")

		if clientToken == "" {
			r := models.Result{
				Status:  "Authentication failure",
				Message: "No Authorization header provided 😰😰😰",
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(r)
			return
		}

		if !strings.HasPrefix(clientToken, "Bearer") {
			r := models.Result{
				Status:  "Authentication failure",
				Message: "Add 'Bearer' to the token 👍🏾👍🏾👍🏾",
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(r)
			return
		}

		if strings.HasPrefix(clientToken, "Bearer") {

			tk := strings.TrimPrefix(clientToken, "Bearer ")

			claims, err := helpers.ValidateToken(tk)

			if err != "" {
				r := models.Result{
					Status:  "Authentication failure",
					Message: err,
				}
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(r)
				return
			}

			r.Header.Set("email", claims.Email)
			r.Header.Set("first_name", claims.FirstName)
			r.Header.Set("last_name", claims.LastName)
			r.Header.Set("uid", claims.UID)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)

		}

	})
}
