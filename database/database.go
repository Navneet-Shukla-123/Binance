package database

import (
	"binance/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// ConnectToDB function will connect to database
func ConnectToDB() {

	username := "posttest"
	password := "test1234"
	host := "localhost"
	port := "5432"
	databaseName := "Binance"

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		username, password, host, port, databaseName)

	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database!")

	sql := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
        password TEXT NOT NULL
    );`

	_, err = db.Exec(sql)
	if err != nil {
		log.Println("Error in Creating the table ", err)
		return
	}

	log.Println("Table created successfully!!")
}

// InsertUserToDB function will insert the user to DB
func InsertUserToDB(email, password string) error {
	sql := `INSERT INTO users ("email","password")
    VALUES($1,$2)`

	_, err := db.Exec(sql, email, password)
	if err != nil {
		return err
	}
	return nil
}

// GetUser function will get the total user count by email
func GetUser(email string) (models.User, error) {
	sql := `Select email,password from users where email=$1`
	var user models.User

	err := db.QueryRow(sql, email).Scan(&user.Email, &user.Password)
	if err != nil {

		return user, err
	}

	return user, nil
}

// InsertTokenToDB function will insert token of user to tokens table
func InsertTokenToDB(email, token string, expiry time.Time) error {
	sql := `INSERT into tokens("email","tokens","expires_at")
	      VALUES($1,$2,$3);`

	_, err := db.Exec(sql, email, token, expiry)
	if err != nil {
		return err
	}
	return nil
}

// GetTokenDatabyToken function get the token data by token
func GetTokenDatabyToken(token string) (models.Tokendata, error) {
	sql := `Select email,tokens,expires_at from tokens where tokens=$1;`

	var data models.Tokendata

	err := db.QueryRow(sql, token).Scan(&data.Email, &data.Token, &data.Expires_At)
	if err != nil {

		return data, err
	}

	return data, nil

}

// UpdateToken will update the token
func UpdateToken(email, token string, expiry time.Time) error {
	sql := `UPDATE tokens set tokens=$1,expires_at=$2 where email=$3;`
	_, err := db.Exec(sql, token, expiry, email)
	if err != nil {
		return err
	}
	return nil
}

// GetTokenCount function will used to check  if the tokken is present
func GetTokenCount(email string) (int, error) {
	var count int

	sql := `SELECT count(email) from tokens where email=$1;`

	err := db.QueryRow(sql, email).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
