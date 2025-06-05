# **DNS Server**

A simple DNS server built with Go, designed for easy deployment via Docker.  

## **Project Overview**
This project provides a **local DNS server** using Go and YAML configuration. It’s containerized with **Docker Compose** for streamlined setup.

## **Features**
✔ Custom **hostname-to-IP mappings** via YAML  
✔ Lightweight Go-based DNS resolution  
✔ Fully containerized for easy deployment  
✔ Configurable via Docker Compose  

⚠️ **Warning:** You should not use this project for production purposes. It is intended for learning and development purposes only.

---

## **Project Structure**
```
dns-server/
├── dns_records.yml      # YAML file containing DNS records
├── docker-compose.yml   # Docker Compose setup
├── Dockerfile           # Docker build instructions
├── go.mod               # Go module dependencies
├── go.sum               # Go package checksums
├── LICENSE              # Project license
├── main.go              # Main DNS server code (Go)
├── README.md            # Project documentation
```

---

## **Getting Started**
### **1️⃣ Clone the Repository**
```bash
git clone https://github.com/petrusjohannesmaas/dns-server.git
cd dns-server
```

### **2️⃣ Build the Docker Image**
```bash
docker build -t dns-server .
```

### **3️⃣ Run with Docker Compose**
```bash
docker compose up -d
```

### **4️⃣ Test DNS Resolution**
Use `dig` or `nslookup` to verify DNS functionality:
```bash
dig @localhost dev-machine.local
```
```bash
nslookup dev-machine.local localhost
```

---

## **Configuration**
Modify `dns_records.yml` to update hostname mappings:
```yaml
records:
  - hostname: "dev-machine.local"
    ip: "192.168.0.100"
  - hostname: "server.local"
    ip: "192.168.0.200"
```

---

## **License**
This project is licensed under the **GNU GENERAL PUBLIC LICENSE Version 3**. See the [LICENSE](LICENSE) file for details.

## **Contributing**
Feel free to **open an issue** or submit a **pull request** to improve the project!

---

### Future Enhancements

✔ **Troubleshooting**: Include troubleshooting steps in the README.
✔ **Front end configuration**: Add a web interface for easy hostname management.
✔ **Automatic reloads**: Ensure new records apply **without restarting** the server.
✔ **Logging**: Capture query analytics.
✔ **Security**: Add authentication for managing DNS entries.
✔ **Persistent storage**: Mount `dns_records.yaml` so records survive container restarts.
✔ **Caching**: Implement a caching mechanism to improve performance.