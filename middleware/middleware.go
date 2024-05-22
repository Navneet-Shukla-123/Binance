package middleware

import (
	"binance/database"
	"errors"
	"log"
	"net/http"
	"time"
)

// Authenticate middleware to authenticate the user
func AuthenticateByToken(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("Middleware is running ... ")

		headerData := r.Header.Get("token")
		if len(headerData) == 0 {
			w.Write([]byte("Token is not available"))
			return

		}

		data, err := database.GetTokenDatabyToken(headerData)
		if err != nil {
			log.Println("Error in getting the token data ", err)
			w.Write([]byte("Error  in getting the token data"))
			return
		}

		log.Println("Expiry time is ", data.Expires_At)

		log.Println("Current time is ", time.Now())

		currentTime := time.Now()
		expiryTime := data.Expires_At

		if expiryTime.Before(currentTime) {
			log.Println("Token has expired")
			w.Write([]byte("Token has expired.Please login"))
			return
		}

		f(w, r)
	}
}

// AuthenticateByCookie middleware will authenticate the user by email
func AuthenticateByCookie(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("cookie")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Println("Cookie is not present")
				w.Write([]byte("Please login again"))
				return

			}
			log.Println("Error in getting the cookie ", err)
			w.Write([]byte("Error in getting the cookie"))
			return
		}

		log.Println("Cookie is ", cookie)
		f(w, r)
	}
}
