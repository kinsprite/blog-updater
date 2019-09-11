package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

var serverSignature string
var shellScriptFile string

func init() {
	serverSignature = os.Getenv("SERVER_SIGNATURE")
	shellScriptFile = os.Getenv("SHELL_SCRIPT_FILE")
}

func main() {
	addr := os.Getenv("LISTENING_ADDRESS")

	if addr != "" {
		addr = ":8080"
	}

	mainRouter().Run(addr)
}

func mainRouter() *gin.Engine {
	engine := gin.Default()

	engine.POST("/github-webhooks", handlerGithubWebhooks)

	return engine
}

func handlerGithubWebhooks(c *gin.Context) {
	req := c.Request
	event := req.Header.Get("X-GitHub-Event")
	signature := req.Header.Get("X-Hub-Signature")

	if signature != serverSignature {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	if event == "ping" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Pong",
		})
	} else if event == "push" {
		doUpdate()
		c.JSON(http.StatusOK, gin.H{
			"message": "Done",
		})
	}
}

func doUpdate() {
	out, err := exec.Command(shellScriptFile).Output()

	if err != nil {
		log.Println("[ERROR]  ", err)
	} else {
		log.Println("[INFO]  ", out)
	}
}
