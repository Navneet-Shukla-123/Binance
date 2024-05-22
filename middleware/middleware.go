package middleware

import (
	"binance/database"
	"log"
	"net/http"
	"time"
)

// Authenticate middleware to authenticate the user
func Authenticate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
