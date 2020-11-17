package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Logger for the app
var Logger *log.Logger

func main() {
	// open file for debugging
	logFile, err := os.OpenFile("/tmp/goScrapeAnsibleMetrics.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening file: %v\n", err)
	}
	defer logFile.Close()
	Logger = log.New(logFile, "", log.LstdFlags)

	if len(os.Args) < 2 {
		Logger.Println(fmt.Sprint(os.Args))
		Logger.Println("You must pass in an API token when calling this script.")
		Logger.Println("Usage: ./goScrapeAnsibleMetrics 1234asdf1234")
		os.Exit(0)
	}

	rawMetrics, err := getMetrics(os.Args[1])
	if err != nil {
		Logger.Println(fmt.Sprintf("There was an error scraping Ansible. Error: %v", err.Error()))
		os.Exit(0)
	}

	Logger.Println(fmt.Sprintf("Received metrics:\n%v", rawMetrics))
	convertMetricsToILP(rawMetrics)
}

func getMetrics(apiToken string) (string, error) {
	// Build the request object
	req, err := http.NewRequest("GET", "https://localhost/api/v2/metrics/", nil)
	if err != nil {
		return "", err
	}

	// Add the API token
	apiTokenField := fmt.Sprintf("Bearer %v", apiToken)
	req.Header.Add("Authorization", apiTokenField)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	// Perform the request
	r, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Check the status code
	if r.StatusCode != 200 {
		return "", fmt.Errorf("Invalid status code from Ansible Tower: %v. ", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", fmt.Errorf("Couldn't read the body of the request: %v", err)
	}

	return string(b), nil
}

func convertMetricsToILP(rawMetrics string) {
	metrics := strings.Split(rawMetrics, "\n")

	for _, metric := range metrics {
		if len(metric) > 1 {
			if metric[0] != '#' {
				noQuotes := strings.ReplaceAll(metric, "\"", "")
				cleanMetric := strings.ReplaceAll(noQuotes, "{", ",")
				newMetric := strings.ReplaceAll(cleanMetric, "}", "")
				newestMetric := strings.ReplaceAll(newMetric, " ", " value=")
				unix := time.Now().Unix()

				final := fmt.Sprintf("%v %v", newestMetric, unix)
				fmt.Println(final)

				Logger.Println("Printed metric: " + final)
			}
		}
	}
}
