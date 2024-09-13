package main

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "mc_allowlist"
)

type UserForm struct {
    Username string `form:"username" binding:"required"`
    Message string `form:"msg"`
}

type UserJson struct {
    Username string `json:"username"`
}

var db *sql.DB

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

    var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

    r := gin.Default()
    r.SetTrustedProxies(nil)
    r.POST("/newUserRequest", handleNewUserRequest)
    r.POST("/checkUsername", handleUserValidation)
    //uri := fmt.Sprintf("%s:%s",address, port)
    r.Run("localhost:8080")

}

func handleNewUserRequest(c *gin.Context) {
    var newUser UserForm;

    if err := c.ShouldBind(&newUser); err != nil {
        return
    }
    // Revalidate Username
    if res, err := checkValidUser(newUser.Username); res != true || err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to validate user"})
        return
    }

    //check if user exists in DB
    var existingUsername string
    row := db.QueryRow("SELECT username FROM users WHERE username = $1", newUser.Username)
    switch err := row.Scan(&existingUsername); err {
    case sql.ErrNoRows:
        break
    case nil:
       if existingUsername != "" {
           c.JSON(http.StatusConflict, gin.H{"error": "Username already requested"})
           return
       }
    default:
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
        return
    }

    requestDate := time.Now().UTC()
    approved := false;
    fmt.Println(newUser.Username, newUser.Message)
    sqlStatement := `
    INSERT INTO users (username, message, request_date, approval_date, approved)
    VALUES ($1, $2, $3, $4, $5)`
    _, err := db.Exec(sqlStatement, newUser.Username, newUser.Message, requestDate, nil, approved)

    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    // Poke discord bot 

    c.JSON(http.StatusOK, gin.H{"message": "User request submitted - awaiting approval"})
}

func handleUserValidation(c *gin.Context) {
    var newUser UserJson;

    if err := c.BindJSON(&newUser); err != nil {
        return
    }

    res, err := checkValidUser(newUser.Username)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to validate user"})
        return
    }

    if res {
        c.JSON(http.StatusOK, gin.H{"message": "Username is valid"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error" : "Username is invalid"})
    }
}

func checkValidUser(username string) (bool, error) {
    uri := "https://api.mojang.com/users/profiles/minecraft/" + username

    resp, err := http.Get(uri)
    
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        return true, nil
    }
    return false, nil
}
