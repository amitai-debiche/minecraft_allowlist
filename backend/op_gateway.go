package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
)

type UserForm struct {
	Username string `form:"user_username" binding:"required"`
	Message  string `form:"user_message"`
}

type UserJson struct {
	Username string `json:"username"`
}
type UserApproval struct {
	Username string `json:"username"`
	Status   bool   `json:"status"`
}

var db *sql.DB

func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
        expectedApiKEY := "Bearer " + os.Getenv("API_KEY")

        if apiKey != expectedApiKEY {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
        }

		c.Next()
	}
}

func handleNewUserRequest(c *gin.Context) {
	var newUser UserForm

	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading form"})
		return
	}

	// Revalidate Username
	if res, err := checkValidUser(newUser.Username); res != true || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate user"})
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
	approved := false
	fmt.Println(newUser.Username, newUser.Message)
	sqlStatement := `
    INSERT INTO users (username, message, request_date, approval_date, approved)
    VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, newUser.Username, newUser.Message, requestDate, nil, approved)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
    fmt.Printf("added to sql\n")

	// Poke discord bot
	apiKey := os.Getenv("API_KEY")
    url := "http://127.0.0.1:5000/send_message"

    data := map[string]string{
        "username" : newUser.Username,
        "message" : newUser.Message,
    }
    jsonData, _ := json.Marshal(data)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}

	// Set content type and authorization headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message to Discord bot"})
		return
	}

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
		fmt.Printf("received non-OK status: %s", resp.Status)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Discord bot responded with error"})
        return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User request submitted - awaiting approval"})
}

func handleUserApproval(c *gin.Context) {
	var userStatus UserApproval

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
    fmt.Printf("approved: %t", approved)
	if approved {
		if err := sendAllowlistRequest(userStatus.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to allowlist"})
			return
		} else {
			msg := fmt.Sprintf("User %s has been added to the allowlist", userStatus.Username)
			c.JSON(http.StatusOK, gin.H{"message": msg})
			return
		}
	} else {
		msg := fmt.Sprintf("User %s has not been added to the allowlist", userStatus.Username)
		c.JSON(http.StatusNotImplemented, gin.H{"message": msg})
		return
	}
}

func handleUserValidation(c *gin.Context) {
	var newUser UserJson

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	res, err := checkValidUser(newUser.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate user"})
		return
	}

	if res {
		c.JSON(http.StatusOK, gin.H{"message": "Username is valid"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is invalid"})
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

func sendAllowlistRequest(username string) error {
	url := "http://localhost:1025/addUser"

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

	if resp.StatusCode != http.StatusOK  && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("received non-OK status: %s", resp.Status)
	}

	return nil
}

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
    release_mode := os.Getenv("RELEASE_MODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

    var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Unable to connect to the database")
		panic(err)
	}

    if strings.ToLower(release_mode) == "true" {
        gin.SetMode(gin.ReleaseMode)
    }

	r := gin.Default()
	r.SetTrustedProxies(nil)
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"POST"},
        AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders: []string{"Content-Length"}, 
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

	r.POST("/newUserRequest", apiKeyMiddleware(), handleNewUserRequest)
	r.POST("/checkUsername", apiKeyMiddleware(), handleUserValidation)
	r.POST("/approveUsername", apiKeyMiddleware(), handleUserApproval)
	r.Run("0.0.0.0:8080")

}
