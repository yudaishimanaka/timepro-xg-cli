package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

type Config struct {
	Conf TimeproXgConfig `toml:"timepro-xg"`
}

type TimeproXgConfig struct {
	RequestUrl string `toml:"request_url"`
	UserId     string `toml:"user_id"`
	Password   string `toml:"password"`
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

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
				attendanceRequest(conf, "PUNCH1", "PUNCH1")
				if askForConfirmation() {
					fmt.Printf("今日も一日がんばるぞい!\n")
				}

				return nil
			},
		},
		{
			Name:  "out",
			Usage: "timepro-xg out <- Leaving work!",
			Action: func(c *cli.Context) error {
				fmt.Printf("Do you really leaving work? [y/n] :")
				attendanceRequest(conf, "PUNCH2", "PUNCH2")
				if askForConfirmation() {
					fmt.Printf("今日も一日お疲れ様でした!\n")
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func attendanceRequest(conf Config, pageStatus, process string) error {
	values := url.Values{}
	values.Add("PAGESTATUS", pageStatus)
	values.Add("LoginID", conf.Conf.UserId)
	values.Add("PassWord", conf.Conf.Password)
	values.Add("PROCESS", process)

	req, err := http.NewRequest("POST", conf.Conf.RequestUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("Accept-Language", `ja,en-US;q=0.7,en;q=0.3`)
	req.Header.Add("Referer", conf.Conf.RequestUrl)
	req.Header.Add("Content-Type", `application/x-www-form-urlencoded`)
	req.Header.Add("Connection", `keep-alive`)
	req.Header.Add("Upgrade-Insecure-Request", `1`)

	req.URL.RawQuery = values.Encode()
	_, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		return err2
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
