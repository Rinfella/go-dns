package main

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

// RecordType represents a DNS record type
type RecordType string

// DNS Record types
const (
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	CNAME RecordType = "CNAME"
	MX    RecordType = "MX"
	TXT   RecordType = "TXT"
	NS    RecordType = "NS"
)

type DNSResult struct {
	Domain    string
	Type      RecordType
	Records   []string
	Error     error
	QueryTime time.Duration
}

// Lookup DNS performs a DNS lookup for the specified domain and record type
func LookupDNS(domain string, recordType RecordType, server string) DNSResult {
	result := DNSResult{
		Domain: domain,
		Type:   recordType,
	}

	if server == "" {
		server = "8.8.8.8:53" // Defaults to Google's DNS
	}

	// Create a new message with the appropriate question
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), recordTypeToInt(recordType))
	m.RecursionDesired = true

	// Create a new Client
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	// Send the query and measure the time
	r, rtt, err := c.Exchange(m, server)
	result.QueryTime = rtt

	if err != nil {
		result.Error = fmt.Errorf("DNS query failed: %v", err)
		return result
	}

	// Check if we got a response
	if r == nil {
		result.Error = fmt.Errorf("No response from the DNS server")
		return result
	}

	// Parse the answer section
	for _, ans := range r.Answer {
		switch v := ans.(type) {
		case *dns.A:
			result.Records = append(result.Records, v.A.String())
		case *dns.AAAA:
			result.Records = append(result.Records, v.AAAA.String())
		case *dns.CNAME:
			result.Records = append(result.Records, v.Target)
		case *dns.MX:
			result.Records = append(result.Records, fmt.Sprintf("%d %s", v.Preference, v.Mx))
		case *dns.TXT:
			result.Records = append(result.Records, v.Txt...)
		case *dns.NS:
			result.Records = append(result.Records, v.Ns)
		}
	}
	return result
}

// recordTypeToInt converts a recordType to its corresponding dns.Type constant
func recordTypeToInt(rt RecordType) uint16 {
	switch rt {
	case A:
		return dns.TypeA
	case AAAA:
		return dns.TypeAAAA
	case CNAME:
		return dns.TypeCNAME
	case MX:
		return dns.TypeMX
	case TXT:
		return dns.TypeTXT
	case NS:
		return dns.TypeNS
	default:
		return dns.TypeA
	}
}

// GetSystemDNSServers returns the configured DNS servers on the system
func GetSystemDNSServers() []string {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return []string{"8.8.8.8:53", "1.1.1.1:53"} // Return default servers if we can't read the config
	}

	servers := make([]string, len(config.Servers))
	for _, server := range config.Servers {
		servers = append(servers, net.JoinHostPort(server, "53"))
	}

	if len(servers) == 0 {
		return []string{"8.8.8.8:53", "1.1.1.1:53"} // Return default servers if no servers were found
	}
	return servers
}
