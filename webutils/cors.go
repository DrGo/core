package webutils

import "net/http"

// CORSHandled sets CORS headers and returns true if it handled the (OPTIONS) request
func CORSHandled(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	return r.Method == "OPTIONS"
}

// BasicAuth checks that the request contains a valid username and pwd
// src: https://www.alexedwards.net/blog/basic-authentication-in-go
func BasicAuthUserName(w http.ResponseWriter, r *http.Request) string {
	name, _, ok := r.BasicAuth()
	if ok {
		return name
		// // Calculate SHA-256 hashes for the provided and expected
		// // usernames and passwords.
		// usernameHash := sha256.Sum256([]byte(username))
		// passwordHash := sha256.Sum256([]byte(password))
		// expectedUsernameHash := sha256.Sum256([]byte("your expected username"))
		// expectedPasswordHash := sha256.Sum256([]byte("your expected password"))

		// // Use the subtle.ConstantTimeCompare() function to check if
		// // the provided username and password hashes equal the
		// // expected username and password hashes. ConstantTimeCompare
		// // will return 1 if the values are equal, or 0 otherwise.
		// // Importantly, we should to do the work to evaluate both the
		// // username and password before checking the return values to
		// // avoid leaking information.
		// usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		// passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expe
	}
	// If the Authentication header is not present, is invalid, or the
	// username or password is wrong, then set a WWW-Authenticate
	// header to inform the client that we expect them to use basic
	// authentication and send a 401 Unauthorized response.
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return ""
}
