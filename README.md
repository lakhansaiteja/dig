# Dig 

Dig app provides DNS lookup for domains and IP addresses

## APIs
### DNS Records
Returns DNS records of for the combination of domain, DNS Type (ex: MX, AAAA, PTR, NS) and server
* **URL**

  `localhost:8080/api/v1/dig`

* **Method:**

  `POST`
* **Example:**
```go
// example request
{
    "domain": "google.com",
    "type": "NS",
    "server": "1.1.1.1"
}

// example response
{
    "Records": [
        "ns4.google.com.",
        "ns1.google.com.",
        "ns3.google.com.",
        "ns2.google.com."
    ]
}
```
### History
Returns history of previous reuests upto 30 records with timestamps
* **URL**

  `localhost:8080/api/v1/dig/history`

* **Method:**

  `GET`
* **Example:**
```go
// example response
[
    {
    "Request": {
        "domain": "google.com",
        "type": "NS",
        "server": "1.1.1.1"
        },
    "Timestamp": "2023-04-13 01:01:09.835604 -0700 PDT m=+210.057975543"
    },
    {
    "Request": {
        "domain": "google.com",
        "type": "AAAA",
        "server": "1.1.1.1"
        },
    "Timestamp": "2023-04-13 00:57:40.877426 -0700 PDT m=+1.090100876"
    }
]
```

## Deployment

`make docker` to run the app in docker  
`make build` to build the binary for different OS and architectures  
`make k8s` to deploy the app in kubernetes


## Built With

* [Go](https://golang.org/) - The programming language used


## Author

* **Lakhan Kamireddy** - [LinkedIn](https://www.linkedin.com/in/lakhansaiteja/)
