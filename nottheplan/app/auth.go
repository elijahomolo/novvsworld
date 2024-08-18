package app

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Claims   *Claims
	Writer   http.ResponseWriter
}

func (a *Auth) createSession() {
	// create a user session
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(10 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	a.Claims = &Claims{
		Username: a.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		auth.Writer.WriteHeader(http.StatusInternalServerError)
		errorPage(auth.Writer, http.StatusInternalServerError)
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(auth.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func (a *Auth) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var tknStr string
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			errorPage(w, http.StatusUnauthorized)
		}
		log.Printf("Failed to get token: %v", err)
		errorPage(w, http.StatusBadRequest)
	}
	// Get the JWT string from the cookie
	tknStr = c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			errorPage(w, http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed to parse token: %v", err)
		errorPage(w, http.StatusBadRequest)
	}
	if !tkn.Valid {
		errorPage(w, http.StatusUnauthorized)
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	//write, err := w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	//if err != nil {
	//	return nil, nil
	//}
	//return claims, http.StatusOK, nil
}
