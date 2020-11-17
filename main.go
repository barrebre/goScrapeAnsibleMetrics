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
		text := []byte("couldn't query endpoint")
		err := ioutil.WriteFile("/tmp/goScrapeAnsibleMetricsErr", text, 0644)
		if err != nil {
			fmt.Println("Coudln't write to file")
		}
		return "", err
	}

	// Check the status code
	if r.StatusCode != 200 {
		text := []byte(fmt.Sprintf("Invalid status code from Ansible Tower: %v. ", r.StatusCode))
		err := ioutil.WriteFile("/tmp/goScrapeAnsibleMetricsErr", text, 0644)
		if err != nil {
			fmt.Println("Coudln't write to file")
		}
		return "", fmt.Errorf("Invalid status code from Ansible Tower: %v. ", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		text := []byte("Couldn't read the body of the request")
		err := ioutil.WriteFile("/tmp/goScrapeAnsibleMetricsErr", text, 0644)
		if err != nil {
			fmt.Println("Coudln't write to file")
		}
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

				final := fmt.Sprintf("%v %v\n", newestMetric, unix)

				fmt.Println(final)

				text := []byte(fmt.Sprintln(final))
				_ = ioutil.WriteFile("/tmp/goScrapeAnsibleMetrics", text, 0644)
			}
		}
	}

	text := []byte("Finished writing all\n\n")
	err := ioutil.WriteFile("/tmp/goScrapeAnsibleMetricsOut", text, 0644)
	if err != nil {
		fmt.Println("Coudln't write to file")
	}
}
