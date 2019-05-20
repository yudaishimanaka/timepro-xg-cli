package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli"
)

var (
	requestUrl string
	userId     string
	password   string
)

func main() {
	requestUrl = os.Getenv("TIMEPROXG_REQUEST_URL")
	userId = os.Getenv("TIMEPROXG_USERID")
	password = os.Getenv("TIMEPROXG_PASSWORD")

	app := cli.NewApp()
	app.Name = "timepro-xg"
	app.Usage = "TimePro-XG CLI Tool"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "in",
			Usage: "timepro-xg in <- Going to work!",
			Action: func(c *cli.Context) error {
				fmt.Printf("Do you really go to work? [y/n] :")
				if askForConfirmation() {
					fmt.Printf("今日も一日がんばるぞい!\n")
					if err := attendanceRequest("PUNCH1", "PUNCH1"); err != nil {
						log.Fatal(err)
					}
				}

				return nil
			},
		},
		{
			Name:  "out",
			Usage: "timepro-xg out <- Leaving work!",
			Action: func(c *cli.Context) error {
				fmt.Printf("Do you really leaving work? [y/n] :")
				if askForConfirmation() {
					fmt.Printf("今日も一日お疲れ様でした!\n")
					if err := attendanceRequest("PUNCH2", "PUNCH2"); err != nil {
						log.Fatal(err)
					}
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func attendanceRequest(pageStatus, process string) error {
	values := url.Values{}
	values.Add("PAGESTATUS", pageStatus)
	values.Add("PROCESS", process)
	values.Add("LoginID", userId)
	values.Add("PassWord", password)

	_, err := http.PostForm(requestUrl, values)
	if err != nil {
		return err
	}

	return nil
}

// Awesome Code, thank you!!
// https://gist.github.com/albrow/5882501
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if response == "" {
		fmt.Printf("Please type yes or no and then press enter: ")
		return askForConfirmation()
	}
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Printf("Please type yes or no and then press enter: ")
		return askForConfirmation()
	}
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
