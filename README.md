# DNS Server

---

## **Step 1: Build the DNS Server in Go**
### **1. Install Dependencies**
Ensure you have Go installed, then install the DNS and YAML libraries:

```bash
go get github.com/miekg/dns gopkg.in/yaml.v3
```

### **2. Define the YAML File for DNS Records**
Create `dns_records.yaml` for hostname-IP mappings:

```yaml
records:
  - hostname: "dev-machine.local"
    ip: "192.168.1.100"
  - hostname: "server.local"
    ip: "192.168.1.200"
```

### **3. Implement the DNS Server**
Create `dns_server.go`:

```go
package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"github.com/miekg/dns"
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
	data, err := ioutil.ReadFile("dns_records.yaml")
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
```

### **4. Run Your DNS Server**
Start it:
```bash
go run dns_server.go
```

Test it:
```bash
dig @localhost dev-machine.local
```

---

## **Step 2: Build the Web Front End**
### **1. Create a Simple HTML Interface**
Make `index.html`:

```html
<form id="dnsForm">
    <label>Hostname:</label>
    <input type="text" id="hostname" required>

    <label>IP Address:</label>
    <input type="text" id="ipAddress" required>

    <button type="submit">Add Record</button>
</form>

<script>
document.getElementById("dnsForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = {
        hostname: document.getElementById("hostname").value,
        ip: document.getElementById("ipAddress").value,
    };

    const response = await fetch("/api/add-dns-record", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });

    if (response.ok) {
        alert("DNS record added!");
    }
});
</script>
```

### **2. Build an API to Update YAML**
Modify `dns_server.go`:

```go
import (
	"net/http"
	"encoding/json"
)

func saveRecords() {
	data, err := yaml.Marshal(Config{Records: convertMapToSlice(dnsRecords)})
	if err != nil {
		fmt.Println("Error encoding YAML:", err)
		return
	}

	err = ioutil.WriteFile("dns_records.yaml", data, 0644)
	if err != nil {
		fmt.Println("Error writing YAML file:", err)
	}
}

func addDNSRecord(w http.ResponseWriter, r *http.Request) {
	var record DNSRecord
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	dnsRecords[record.Hostname+"."] = record.IP
	saveRecords() // Persist changes to YAML

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/api/add-dns-record", addDNSRecord)
	fmt.Println("DNS API running on port 8080")
	http.ListenAndServe(":8080", nil)
}
```

### **3. Run the Web Server**
Start it:
```bash
go run dns_server.go
```

Test via curl:
```bash
curl -X POST http://localhost:8080/api/add-dns-record \
     -H "Content-Type: application/json" \
     -d '{"hostname":"test.local","ip":"192.168.1.150"}'
```

---

## **Step 3: Integrate with Your Tenda Router**
### **1. Configure DNS Settings**
- Log into the router (`192.168.0.1`).
- Go to **Advanced Settings** → **DNS Settings**.
- In **Preferred DNS Server**, enter your DNS server’s IP (`192.168.1.100`).
- In **Alternate DNS Server**, set a fallback (`8.8.8.8` for Google DNS).
- Save and apply.

### **2. Restart the Router**
Reboot to ensure DNS settings are active.

### **3. Test Network-Wide DNS Resolution**
Run:
```bash
dig @192.168.1.100 dev-machine.local
```

---

## **Next Enhancements**
✔ **Automatic reloads**: Ensure new records apply **without restarting** the server.  
✔ **Persistent database storage**: Store records in **SQLite** instead of just YAML.  
✔ **Security measures**: Add authentication to prevent unauthorized changes.  

Containerizing your **Go-based DNS server** with Docker will make deployment easier across different environments. Here’s how to do it step by step:

---

## **Step 1: Create a Dockerfile**
Inside your project folder, create a `Dockerfile`:

```dockerfile
# Use a lightweight Go image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the Go module files & download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o dns-server dns_server.go

# Expose DNS port (UDP)
EXPOSE 53/udp

# Expose HTTP API port (optional)
EXPOSE 8080

# Start the DNS server
CMD ["./dns-server"]
```

---

## **Step 2: Create a Docker Compose File (Optional)**
If you want to manage multiple services (e.g., a database or a proxy), use `docker-compose.yaml`:

```yaml
version: '3.8'

services:
  dns-server:
    build: .
    container_name: dns_server
    restart: always
    ports:
      - "53:53/udp"
      - "8080:8080"
    volumes:
      - ./dns_records.yaml:/app/dns_records.yaml
```

---

## **Step 3: Build & Run the Container**
Run these commands in your project directory:

```bash
docker build -t go-dns-server .
docker run -d --name dns-server -p 53:53/udp -p 8080:8080 go-dns-server
```

If using **Docker Compose**, start everything with:
```bash
docker-compose up -d
```

---

## **Step 4: Verify the Container is Running**
Check logs:
```bash
docker logs dns-server
```

Test DNS resolution:
```bash
dig @localhost dev-machine.local
```

---

## **Step 5: Set Up Your Router**
- Go to your **Tenda router** admin panel.
- Locate **DNS settings**.
- Set the **Preferred DNS Server** to your **container’s IP** (e.g., `192.168.1.100`).
- Apply & restart the router.

---

## **Next Enhancements**
✔ **Persistent storage**: Mount `dns_records.yaml` so records survive container restarts.  
✔ **Secure API**: Add authentication for managing DNS entries.  
✔ **Logging**: Capture query analytics. 
