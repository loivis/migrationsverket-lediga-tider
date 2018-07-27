package main

import (
	"fmt"
	"log"

	helpers "github.com/loivis/go-helpers"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func startChrome() selenium.WebDriver {
	var err error
	opts := []selenium.ServiceOption{
	// selenium.StartFrameBuffer(),
	// selenium.Output(os.Stderr),
	}
	// selenium.SetDebug(true)

	port := helpers.RandomIntBetween(9000, 9999)
	_, err = selenium.NewChromeDriverService("/usr/local/bin/chromedriver", port, opts...)
	if err != nil {
		log.Panicf("Error starting the ChromeDriver server: %v", err)
	}
	// defer service.Stop()

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	// disable image loading
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
		},
	}
	caps.AddChrome(chromeCaps)

	// create remote client
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	// defer webDriver.Quit()

	// add custom cookies
	// webDriver.AddCookie(&selenium.Cookie{
	// 	Name:  "name",
	// 	Value: "value",
	// })

	return webDriver
}
