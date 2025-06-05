package main

import (
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

type DNSRecord struct {
	Hostname string `yaml:"hostname"`
	IP       string `yaml:"ip"`
}

type Config struct {
	Records []DNSRecord `yaml:"records"`
}

var dnsRecords map[string]string

func loadRecords() {
	data, err := os.ReadFile("dns_records.yml")
	if err != nil {
		fmt.Println("Error reading YAML:", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return
	}

	dnsRecords = make(map[string]string)
	for _, record := range config.Records {
		dnsRecords[record.Hostname+"."] = record.IP
		fmt.Printf("Loaded: %s -> %s\n", record.Hostname, record.IP)
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	for _, q := range r.Question {
		if ip, found := dnsRecords[q.Name]; found {
			rr := &dns.A{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(ip),
			}
			m.Answer = append(m.Answer, rr)
		}
	}

	w.WriteMsg(m)
}

func main() {
	loadRecords()

	dns.HandleFunc(".", handleDNSRequest)
	server := &dns.Server{Addr: ":53", Net: "udp"}

	fmt.Println("Starting DNS server on port 53...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
