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

	log.Println("Login method is ", routes.LoginMethod)

	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/login", routes.Login)

	http.HandleFunc("/get", middleware.AuthenticateByToken(routes.GetBinanceData))
	http.HandleFunc("/get", middleware.AuthenticateByCookie(routes.GetBinanceData))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error in running the server ", err)
		return
	}

}
