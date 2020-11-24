This app reads from an Ansible Tower metrics endpoint (`api/v2/metrics`) and scrapes the metrics into Influx Line Protocol.

## Usage
`Usage: ./goScrapeAnsibleMetrics -api-token={} -server-url={}`
* **api-token**: The Ansible Tower token to pull metrics
* **server-url**: The Ansible Tower server to pull metrics from

### Telegraf Example
You can use this script as a Telegraf Input to send to any Telegraf Output
```
[[inputs.exec]]
  commands = [ 
    "/etc/telegraf/telegraf.d/goScrapeAnsible -api-token=<api-token>"
  ]
  data_format = "influx"
```

### Assumptions
* This will print metrics in ILP (to be used by Telegraf)
* Port 443 is expected for the Ansible UI
* You must generate a token which can query the metrics

### Dev Build Flag
To build for Linux systems
`env GOOS=linux GOARCH=386 go build -o goScrapeAnsible main.go`