package main

import (
	"fmt"
	"log"
	"os"
	"time"

	helpers "github.com/loivis/go-helpers"
	"github.com/tebeka/selenium"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var (
	domain, apiKey, pubKey string
	urlFormat              = "https://www.migrationsverket.se/ansokanbokning/valjtyp?sprak=sv&bokningstyp=2&enhet=%s&sokande=3"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	domain = os.Getenv("DOMAIN")
	apiKey = os.Getenv("API_KEY")
	pubKey = os.Getenv("PUBLIC_KEY")
}

func main() {
	browser := startChrome()

	var count int

	for {
		offices := map[string]string{
			"göteborg":   "Z102",
			"sundbybery": "Z209",
		}

		log.Println("count:", count)
		for location, code := range offices {
			log.Println(location, code)

			url := fmt.Sprintf(urlFormat, code)

			err := browser.Get(url)
			if err != nil {
				log.Println(err)
			}

			feedback, err := browser.FindElement(selenium.ByClassName, "feedbackPanelERROR")
			if err == nil {
				msg, _ := feedback.Text()
				log.Println(msg)
				continue
			}

			if location != "göteborg" || (location == "göteborg" && count%3 == 0) {
				sendNotification(location, url)
			}
		}

		sleep()
		count++
	}
}

func sendNotification(location, url string) {
	mail := mailgun.NewMailgun(domain, apiKey, pubKey)

	sender := "lediga.tider@" + domain
	subject := "lediga tider för " + location
	body := url
	recipient := "migrationsverket@" + domain

	message := mail.NewMessage(sender, subject, body, recipient)

	resp, id, err := mail.Send(message)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("message send: ID(%s) Resp(%s)", id, resp)
}

func sleep() {
	n := helpers.RandomIntBetween(3, 7)
	log.Printf("sleep %d minutes for next query", n)
	time.Sleep(time.Duration(n) * time.Minute)
}
