# DNS Server

### **1️⃣ Adding a Configuration File**
Instead of hardcoding DNS records, we’ll use a configuration file (like CoreDNS does) to make your DNS server more flexible.

#### **Step 1: Create a Config File**
Create a `dns_config.json` file:
```json
{
    "records": {
        "desktop.local": "192.168.1.100",
        "laptop.local": "192.168.1.101",
        "printer.local": "192.168.1.102"
    }
}
```

#### **Step 2: Modify Your Go Code to Read the Config**
Update `server.go` to read from the config file:
```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net"
    "github.com/miekg/dns"
)

// Load DNS records from config file
func loadConfig(filename string) (map[string]string, error) {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var config map[string]map[string]string
    err = json.Unmarshal(file, &config)
    return config["records"], err
}

// Handle DNS queries
func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg, records map[string]string) {
    msg := new(dns.Msg)
    msg.SetReply(r)

    for _, q := range r.Question {
        if q.Qtype == dns.TypeA {
            if ip, exists := records[q.Name]; exists {
                msg.Answer = append(msg.Answer, &dns.A{
                    Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 300},
                    A:   net.ParseIP(ip),
                })
            }
        }
    }
    w.WriteMsg(msg)
}

func main() {
    records, err := loadConfig("dns_config.json")
    if err != nil {
        fmt.Println("Failed to load config:", err)
        return
    }

    dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
        handleDNSQuery(w, r, records)
    })

    server := &dns.Server{Addr: ":53", Net: "udp"}
    fmt.Println("Starting DNS server on port 53...")

    if err := server.ListenAndServe(); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}
```
✅ Now your DNS server reads records from `dns_config.json`, making it **easier to update**!

---

### **2️⃣ Containerizing the Application**
Now let’s containerize your Go-based DNS server using **Podman**.

#### **Step 1: Create a Dockerfile**
Create a `Dockerfile`:
```dockerfile
FROM golang:latest

WORKDIR /app

COPY server.go dns_config.json .   # Copy files into the container

RUN go mod init mydns && go mod tidy
RUN go build -o dns-server server.go

CMD ["/app/dns-server"]
```

#### **Step 2: Build the Container**
Run:
```sh
podman build -t mydns .
```

#### **Step 3: Run Your DNS Server Container**
Start the container:
```sh
podman run -d --name mydns-container \
    -p 53:53/udp \
    -v ./dns_config.json:/app/dns_config.json \
    mydns
```

🔹 **Your Go-based DNS server is now running in a Podman container!**  
🔹 **You can modify `dns_config.json` without changing the container itself.**

### Future enhancements:

* Add logging
* Add configuration file syntax check
