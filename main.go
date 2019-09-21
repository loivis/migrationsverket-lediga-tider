package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

var (
	apiKey, domain string

	urlFormat = "https://www.migrationsverket.se/ansokanbokning/valjtyp?sprak=sv&bokningstyp=2&enhet=%s&sokande=3"

	offices = map[string]string{
		"göteborg":   "Z102",
		"sundbybery": "Z209",
	}

	httpClient = http.Client{Timeout: 3 * time.Second}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	apiKey = os.Getenv("MAILGUN_API_KEY")
	domain = os.Getenv("MAILGUN_DOMAIN")
}

func main() {
	for location, code := range offices {
		log.Printf("%s: %s", location, code)

		url := fmt.Sprintf(urlFormat, code)

		b, err := fetchContent(url)
		if err != nil {
			log.Printf("error fetching %s: %v", location, err)
			continue
		}

		if strings.Contains(string(b), "Det finns inte lediga tider") {
			log.Printf("%s: no time available", location)
			continue
		}

		if location != "göteborg" {
			sendNotification(location, url)
		}
	}
}

func fetchContent(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func sendNotification(location, url string) {
	mail := mailgun.NewMailgun(domain, apiKey)

	sender := "lediga.tider@" + domain
	subject := "lediga tider för " + location
	body := url
	recipient := "migrationsverket@" + domain

	message := mail.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, id, err := mail.Send(ctx, message)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("message sent: ID(%s) Resp(%s)", id, resp)
}
