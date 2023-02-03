package main

import (
	"context"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

// to run on every route. Checks if the user is logged in. If so, save to context. Either way, continue to the next handler
func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get cookie
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			ctx := context.WithValue(r.Context(), "mealplanner_current_user", -1)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		// Check remember token validity
		user, err := a.services.Users.ByRemember(cookie.Value)
		if err != nil {
			ctx := context.WithValue(r.Context(), "mealplanner_current_user", -1)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		// Add context if authenticated
		ctx := context.WithValue(r.Context(), "mealplanner_current_user", user.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a *application) authorise(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("mealplanner_current_user").(int)
		if !ok || !(user >= 0) {
			http.Redirect(w, r, "/login", http.StatusFound) // Redirect to login if no context
			return
		}
		next.ServeHTTP(w, r)
	})
}
