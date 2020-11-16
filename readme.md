This app reads a localhost instance of Ansible Tower and scrapes the metrics into ILP.

### Assumptions
* This script must run on the Ansible Tower host
* Port 443 is expected for the Ansible UI
* You must generate a token which can query the metrics

### Dev Build Flag
To build for Linux systems
`env GOOS=linux GOARCH=386 go build -o goScrapeAnsible main.go`