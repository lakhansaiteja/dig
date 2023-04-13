package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"time"
)

const defaultDNS = "8.8.8.8"

type DNSRequest struct {
	Domain string `json:"domain,omitempty"`
	Type   string `json:"type,omitempty"`
	Server string `json:"server,omitempty"`
}

type DNSResponse struct {
	Records []string
	Error   *string `json:"error,omitempty"`
}

// dnsRecords is a POST API handler that handles dns requests
func dnsRecords(c *gin.Context) {
	req := DNSRequest{}
	err := c.BindJSON(&req)
	if err != nil || req.Type == "" || req.Domain == "" { //TODO is the nil check required
		c.JSON(http.StatusBadRequest, "invalid request")
	}

	historyRecord := HistoryRecord{Request: req, Timestamp: time.Now().String()}

	mu.Lock() // locking to prevent race condition
	dll.PushBack(historyRecord)
	for dll.Len() > 30 {
		dll.Remove(dll.Front())
	}
	mu.Unlock()

	// using default DNS as fallback when no DNS is provided
	if req.Server == "" {
		req.Server = defaultDNS
	}

	IPvalues, err := lookupRecords(req)
	res := DNSResponse{
		Records: IPvalues,
	}

	var status int
	if err != nil {
		status = http.StatusServiceUnavailable
		errVal := err.Error()
		res.Error = &errVal
	} else {
		status = http.StatusOK
	}

	c.JSON(status, res)
}

// lookupRecords looks up records for the provided domain/ip, DNS server and Type
func lookupRecords(req DNSRequest) ([]string, error) {
	server := req.Server
	i := strings.Index(req.Server, ":")
	if i < 0 { //no colon
		server = fmt.Sprintf("%s:53", req.Server) //using port 53 for DNS by default
	}

	customResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * time.Duration(10),
			}
			return d.DialContext(ctx, network, server)
		},
	}

	var IPvalues []string
	switch req.Type {
	case "A":
		IPs, err := customResolver.LookupIP(context.TODO(), "ip4", req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP.String())
		}

		return IPvalues, nil

	case "AAAA":
		IPs, err := customResolver.LookupIP(context.TODO(), "ip6", req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP.String())
		}

		return IPvalues, nil

	case "CNAME":
		cname, err := customResolver.LookupCNAME(context.TODO(), req.Domain)
		if err != nil {
			return IPvalues, err
		}

		IPvalues = append(IPvalues, cname)
		return IPvalues, nil

	case "PTR":
		IPs, err := customResolver.LookupAddr(context.TODO(), req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP)
		}
		return IPvalues, nil

	case "NS":
		IPs, err := customResolver.LookupNS(context.TODO(), req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP.Host)
		}
		return IPvalues, nil

	case "MX", "ANY":
		IPs, err := customResolver.LookupMX(context.TODO(), req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP.Host)
		}
		return IPvalues, nil

	case "SRV":
		_, IPs, err := customResolver.LookupSRV(context.TODO(), "", "", req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP.Target)
		}
		return IPvalues, nil

	case "TXT":
		IPs, err := customResolver.LookupTXT(context.TODO(), req.Domain)
		if err != nil {
			return IPvalues, err
		}

		for _, IP := range IPs {
			IPvalues = append(IPvalues, IP)
		}
		return IPvalues, nil
	}

	return IPvalues, errors.New("invalid/unknown DNS type")
}
