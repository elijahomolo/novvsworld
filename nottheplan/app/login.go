package app

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

var auth Auth

// Create the JWT key used to create the signature
// create env variable for jwt key
var jwtKey = []byte("my_secret_key")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	//load the login page
	if r.Method != http.MethodPost {
		pd := struct {
			Success bool
			Message string
		}{
			Success: false,
			Message: "Please log in",
		}
		executeTemplate("templates", "login.html", w, pd)
		return
	}

	credentials := Member{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	err := credentials.VerifyUser(credentials)
	if err != nil {
		log.Printf("Failed to verify user: %v", err)
		errorPage(w, http.StatusInternalServerError)
	}

	auth.Username = credentials.Username
	auth.Writer = w

	auth.createSession()

	pd := struct {
		Success bool
		Message string
	}{
		Success: true,
		Message: "You have been logged in",
	}

	executeTemplate("templates", "login.html", w, pd)
}

//func createSession(w http.ResponseWriter, username string) error {
//	// create a user session
//	// Declare the expiration time of the token
//	// here, we have kept it as 5 minutes
//	expirationTime := time.Now().Add(10 * time.Minute)
//
//	// Create the JWT claims, which includes the username and expiry time
//	claims := &Claims{
//		Username: username,
//		RegisteredClaims: jwt.RegisteredClaims{
//			// In JWT, the expiry time is expressed as unix milliseconds
//			ExpiresAt: jwt.NewNumericDate(expirationTime),
//		},
//	}
//
//	// Declare the token with the algorithm used for signing, and the claims
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	// Create the JWT string
//	tokenString, err := token.SignedString(jwtKey)
//	if err != nil {
//		// If there is an error in creating the JWT return an internal server error
//		w.WriteHeader(http.StatusInternalServerError)
//		return err
//	}
//
//	// Finally, we set the client cookie for "token" as the JWT we just generated
//	// we also set an expiry time which is the same as the token itself
//	http.SetCookie(w, &http.Cookie{
//		Name:    "token",
//		Value:   tokenString,
//		Expires: expirationTime,
//	})
//
//	return nil
//}

func Welcome(w http.ResponseWriter, r *http.Request) {
	//tmpl := template.Must(template.ParseFiles("templates/welcome.html"))
	auth.ValidateToken(w, r)

	executeTemplate("templates", "welcome.html", w, &auth.Username)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
