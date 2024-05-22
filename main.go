package main

import (
	"binance/database"
	"binance/middleware"
	"binance/routes"
	"log"
	"net/http"
)

func init() {
	database.ConnectToDB()
}

func main() {

	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/get", middleware.Authenticate(routes.GetBinanceData))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error in running the server ", err)
		return
	}

}
