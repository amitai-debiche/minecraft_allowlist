package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var session_name string

type player struct {
    Username string `json:"username"`
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
        return
    }

    address := os.Getenv("ADDRESS")
    port := os.Getenv("PORT")
    release_mode := os.Getenv("RELEASE_MODE")
    session_name = os.Getenv("TMUX_SESSION_NAME")


    if strings.ToLower(release_mode) == "true" {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()
    r.SetTrustedProxies(nil)
    r.POST("/addUser", addUser)
    uri := fmt.Sprintf("%s:%s",address, port)
    r.Run(uri)
}

func addUser(c *gin.Context){
    var newPlayer player

    if err := c.BindJSON(&newPlayer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request payload",
        })
        return
    }
    if newPlayer.Username == ""{
        return 
    }
    if err := runTmuxCommand(newPlayer.Username); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    message := fmt.Sprintf("Player: %s succesfully added to allowlist", newPlayer.Username)
    c.JSON(http.StatusCreated, gin.H{
        "message": message,
    })
}

func runTmuxCommand(uname string) error {

    command := fmt.Sprintf("/whitelist add %s", uname)
    cmd := exec.Command("tmux", "send-keys", "-t", session_name, command, "C-m")
    fmt.Printf("Executing command: %s\n", cmd.String())

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to run tmux command: %v", err)
    }
    return nil
}




