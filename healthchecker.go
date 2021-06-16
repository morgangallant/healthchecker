// Healthchecker is a simple application to periodically send HTTP requests to an endpoint.
// If the endpoint returns an error, the application is marked as unhealthy and a notification
// is sent to Discord. Healthchecker will continue to issue HTTP requests, and will also send
// a Discord notification once the application is back online. Two required environment variables
// (ENDPOINT, DISCORD_URL), and one optional (SECONDS).
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func must(k string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	panic("missing variable " + k)
}

func fallback(k, f string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return f
}

// Initialized below.
var (
	endpoint   string
	discordUrl string
	seconds    int
)

func init() {
	_ = godotenv.Load()
	endpoint = must("ENDPOINT")
	discordUrl = must("DISCORD_URL")

	// Parse seconds.
	sstr := fallback("SECONDS", "30")
	var err error
	seconds, err = strconv.Atoi(sstr)
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func sendDiscordMessage(msg string) error {
	params := url.Values{}
	params.Set("content", msg)
	resp, err := http.PostForm(discordUrl, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code from discord %d: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func checkUp() error {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code from endpoint %d: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// The current state of the application.
var up = true

func run() error {
	for {
		log.Println("Doing health check.")
		checkErr := checkUp()
		if up && checkErr != nil {
			// The server was previously up, but is now down.
			// 1st fail check again in 2 seconds
			log.Println("First fail")
			time.Sleep(2 * time.Second)
			checkErr = checkUp()
			if up && checkErr != nil {
				// 2nd consecutive fail so we send a message
				log.Println("Failed 2nd health check, notifying.")
				msg := fmt.Sprintf("%s is down... might wanna get on that? (XANDERS TEST)", endpoint)
				if err := sendDiscordMessage(msg); err != nil {
					log.Printf("failed to send discord message: %v", err)
				}
			}
		} else if !up && checkErr == nil {
			// The server was previously down, but is now up.
			log.Println("Back up again, notifying.")
			msg := fmt.Sprintf("%s is back up again, phew!", endpoint)
			if err := sendDiscordMessage(msg); err != nil {
				log.Printf("failed to send discord message: %v", err)
			}
		}
		up = checkErr == nil
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}
