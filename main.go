package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

var listeningAddress string
var serverSecret string
var initScriptFile string
var updateScriptFile string

func init() {
	listeningAddress = os.Getenv("LISTENING_ADDRESS")
	serverSecret = os.Getenv("SERVER_SECRET")
	initScriptFile = os.Getenv("INIT_SCRIPT_FILE")
	updateScriptFile = os.Getenv("UPDATE_SCRIPT_FILE")
}

func parseFlags() {
	listeningAddressFlag := flag.String("LISTENING_ADDRESS", "", "Server listening Address")
	serverSecretFlag := flag.String("SERVER_SECRET", "", "Server authorizing secret")
	initScriptFileFlag := flag.String("INIT_SCRIPT_FILE", "", "Blog init script file on startup")
	updateScriptFileFlag := flag.String("UPDATE_SCRIPT_FILE", "", "Blog updating script file")

	flag.Parse()

	if *listeningAddressFlag != "" {
		listeningAddress = *listeningAddressFlag
	}

	if *serverSecretFlag != "" {
		serverSecret = *serverSecretFlag
	}

	if *initScriptFileFlag != "" {
		initScriptFile = *initScriptFileFlag
	}

	if *updateScriptFileFlag != "" {
		updateScriptFile = *updateScriptFileFlag
	}
}

func main() {
	parseFlags()
	doInit()

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

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// If we have a Secret set, we should check the MAC
	if serverSecret != "" {
		signature := req.Header.Get("X-Hub-Signature")

		if signature == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "403 Forbidden - Missing X-Hub-Signature required for HMAC verification",
			})
			return
		}

		expectedSig := generateSignature(body)

		if !hmac.Equal([]byte(expectedSig), []byte(signature)) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "403 Forbidden - HMAC verification failed",
			})
			return
		}
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

func generateSignature(palyload []byte) string {
	mac := hmac.New(sha1.New, []byte(serverSecret))
	mac.Write(palyload)
	sum := mac.Sum(nil)
	return "sha1=" + hex.EncodeToString(sum)
}

func doInit() {
	runScript(initScriptFile)
}

func doUpdate() {
	runScript(updateScriptFile)
}

func runScript(file string) {
	out, err := exec.Command(file).Output()

	if err != nil {
		log.Println("[ERROR]", err)
	} else {
		log.Println("[INFO]", string(out))
	}
}
