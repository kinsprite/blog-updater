package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

var listeningAddress string
var serverSignature string
var shellScriptFile string

func init() {
	listeningAddress = os.Getenv("LISTENING_ADDRESS")
	serverSignature = os.Getenv("SERVER_SIGNATURE")
	shellScriptFile = os.Getenv("SHELL_SCRIPT_FILE")
}

func parseFlags() {
	listeningAddressFlag := flag.String("LISTENING_ADDRESS", "", "Server listening Address")
	serverSignatureFlag := flag.String("SERVER_SIGNATURE", "", "Server authorizing signature")
	shellScriptFileFlag := flag.String("SHELL_SCRIPT_FILE", "", "Blog updating shell script file")

	flag.Parse()

	if *listeningAddressFlag != "" {
		listeningAddress = *listeningAddressFlag
	}

	if *serverSignatureFlag != "" {
		serverSignature = *serverSignatureFlag
	}

	if *shellScriptFileFlag != "" {
		shellScriptFile = *shellScriptFileFlag
	}
}

func main() {
	parseFlags()

	if listeningAddress == "" {
		listeningAddress = ":8080"
	}

	log.Println("[INFO] Start server on", listeningAddress)
	mainRouter().Run(listeningAddress)
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
		log.Println("[ERROR]", err)
	} else {
		log.Println("[INFO]", string(out))
	}
}
