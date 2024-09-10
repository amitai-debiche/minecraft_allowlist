package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

const TmuxSessionName string = "0";

type player struct {
    Username string `json:"username"`
}

func main() {
    //gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    r.POST("/addUser", addUser)
    r.Run("localhost:8080")
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
        c.IndentedJSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    message := fmt.Sprintf("Player: %s succesfully added to allowlist", newPlayer.Username)
    c.IndentedJSON(http.StatusCreated, gin.H{
        "message": message,
    })
}

func runTmuxCommand(uname string) error {

    command := fmt.Sprintf("/whitelist add %s", uname)
    cmd := exec.Command("tmux", "send-keys", "-t", TmuxSessionName, command, "C-m")

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to run tmux command: %v", err)
    }
    return nil
}




