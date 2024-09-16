package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
type UserApproval struct {
    Username string `json:"username"`
    Status bool `json:"status"`
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
    r.POST("/approveUsername", handleUserApproval)
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

func handleUserApproval(c *gin.Context) {
    var userStatus UserApproval;

    if err := c.BindJSON(&userStatus); err != nil {
        return
    }
    
    //check if user exists in DB
    var existingUsername string
    row := db.QueryRow("SELECT username FROM users WHERE username = $1", userStatus.Username)
    switch err := row.Scan(&existingUsername); err {
    case nil:
        break
    default:
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
        return
    }


    approvalDate := time.Now().UTC()
    approved := userStatus.Status
    sqlStatement := `
    UPDATE users
    SET approval_date = $2, approved = $3
    WHERE username = $1;`

    _, err := db.Exec(sqlStatement, userStatus.Username, approvalDate, approved)
    if err != nil {
        panic(err)
    }

    if approved {
        if err := sendAllowlistRequest(userStatus.Username); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to allowlist"})
            return
        }else {
            msg := fmt.Sprintf("User %s has been added to the allowlist", userStatus.Username)
            c.JSON(http.StatusOK, gin.H{"message": msg});
            return
        }
    }else {
        msg := fmt.Sprintf("User %s has not been added to the allowlist", userStatus.Username)
        c.JSON(http.StatusCreated, gin.H{"message": msg})
        return
    }
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




func sendAllowlistRequest(username string) error{
    url := "http://localhost:8080/addUser"

    payload := map[string]string{"username": username}
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("error marshaling JSON: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-OK status: %s", resp.Status)
    }
        
    return nil
}
