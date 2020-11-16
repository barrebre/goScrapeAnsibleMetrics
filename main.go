package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(os.Args)
		fmt.Println("You must pass in an API token when calling this script.")
		fmt.Println("Usage: ./goScrapeAnsibleMetrics 1234asdf1234")
		os.Exit(0)
	}

	rawMetrics, err := getMetrics(os.Args[1])
	if err != nil {
		fmt.Println("There was an error scraping Ansible. Error: ", err.Error())
		os.Exit(0)
	}

	convertMetricsToILP(rawMetrics)
	// fmt.Println("Received metrics:\n", rawMetrics)
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
				cleanMetric := strings.ReplaceAll(metric, "{", ",")
				new := strings.ReplaceAll(cleanMetric, "}", "")
				final := "ansibleTower:" + new
				println(final)
			}
		}
	}
}
