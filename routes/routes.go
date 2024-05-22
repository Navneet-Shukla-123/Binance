package routes

import (
	"binance/database"
	"binance/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	cookie = iota
	token
)

var LoginMethod int

// Register route will register the user
func Register(w http.ResponseWriter, r *http.Request) {

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println("Error in decoding the request body ", err)
		w.Write([]byte("Error in email  and password"))
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		log.Println("Error in encrypting the password ", err)
		w.Write([]byte("Error in encrypting the password"))
		return
	}
	err = database.InsertUserToDB(u.Email, string(encryptedPassword))
	if err != nil {
		log.Println("Error in saving the user data to database ", err)
		w.Write([]byte("Error in saving the user data to Database"))
		return
	}

	w.Write([]byte("User Registered successfully"))

}

// Login route will login the user
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error in decoding the body ", err)
		w.Write([]byte("Error in decoding the body"))
		return
	}

	userData, err := database.GetUser(user.Email)
	if err != nil {
		log.Println("Error in fetching the user from database ", err)
		w.Write([]byte("Error in fetching the user detail"))
		return
	}
	if userData.Email == "" {
		log.Println("No user for this particular email id")
		w.Write([]byte("No user for this particular email is"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	if err != nil {
		log.Println("Passwords do not match")
		w.Write([]byte("Invalid email or password"))
		return
	}

	// random := uuid.New()
	// expiry_time := time.Now().Add(30 * time.Minute)

	// count, err := database.GetTokenCount(user.Email)
	// if err != nil {
	// 	log.Println("Error in getting the token count for user ", err)
	// 	w.Write([]byte("Error in getting the token count"))
	// 	return
	// }
	// if count != 0 {
	// 	err := database.UpdateToken(user.Email, random.String(), expiry_time)
	// 	if err != nil {
	// 		log.Println("Error in updating the token ", err)
	// 		w.Write([]byte("Error in updating the token"))
	// 		return
	// 	}
	// } else {

	// 	err = database.InsertTokenToDB(user.Email, random.String(), expiry_time)
	// 	if err != nil {
	// 		log.Println("Error in inserting token of the user ", err)
	// 		w.Write([]byte("Error in inserting the token"))
	// 		return
	// 	}
	// }
	random := uuid.New()
	loginMethod := user.LoginMethod

	LoginMethod = loginMethod

	fmt.Println("Login method is ", LoginMethod+10)

	// loginMethod==1 token based, loginMethod==2 cookie based

	if loginMethod == 1 {

		expiry_time := time.Now().Add(30 * time.Minute)

		count, err := database.GetTokenCount(user.Email)
		if err != nil {
			log.Println("Error in getting the token count for user ", err)
			w.Write([]byte("Error in getting the token count"))
			return
		}
		if count != 0 {
			err := database.UpdateToken(user.Email, random.String(), expiry_time)
			if err != nil {
				log.Println("Error in updating the token ", err)
				w.Write([]byte("Error in updating the token"))
				return
			}
		} else {

			err = database.InsertTokenToDB(user.Email, random.String(), expiry_time)
			if err != nil {
				log.Println("Error in inserting token of the user ", err)
				w.Write([]byte("Error in inserting the token"))
				return
			}
		}
		w.Header().Set("token", random.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User login successfull via Token. \n"))
		w.Write([]byte(fmt.Sprintf("Token is %s ", random.String())))
	} else {

		cookie := http.Cookie{
			Name:     "cookie",
			Value:    random.String(),
			Path:     "/",
			MaxAge:   10,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)
		w.Write([]byte("User login successfull via Cookie."))
	}

}

// GetBinanceData function will do the api call to get the data
func GetBinanceData(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("https://api.binance.com/api/v3/depth?symbol=BNBBTC&limit=1000")

	if err != nil {
		log.Println("Error in getting the data from api ", err)
		w.Write([]byte("Error in get request"))
		return
	}

	responseData := models.GetResponse{}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in reading the response body ", err)
		w.Write([]byte("Error in reading the response body"))
		return
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("Error in unmarshalling ", err)
		w.Write([]byte("Error in unmarshalling"))
		return
	}

	jsonData, err := json.MarshalIndent(responseData, "", "\t")
	if err != nil {
		log.Println("Error in converting to json ", err)
		w.Write([]byte("Error in converting to json"))
		return
	}
	headerData := r.Header.Get("token")

	fmt.Println("Cookie is ", cookie)
	fmt.Println("Token is ", token)

	log.Println("Header data is ", headerData)
	w.Write([]byte(jsonData))
}
