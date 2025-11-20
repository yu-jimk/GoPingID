# GoPingID
![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)

A Go-based ping tool that lets you specify a custom ICMP Identifier.

## Features
- Sends ICMP Echo Requests  
- Allows manual specification of the ICMP Identifier (independent from PID)  
- Displays RTT in milliseconds  

## Installation

### 1. Install Go
Make sure Go is installed. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### 2. Clone the repository
```bash
git clone https://github.com/your-username/GoPingID.git
cd GoPingID
```

## Usage
```bash
sudo go run main.go -a <IP_ADDRESS> [options]
```

### Options
| option     | Default      | Description                        | 
| ---------- | :----------- | :-------------------------- | 
| `-a`       | (required)       | Destination IP address           | 
| `-n`       | `3`            | Number of Echo Requests to send        | 
| `-id`      | `PID & 0xffff` | ICMP Identifier (`0–65535`) | 
| `-t`       | `3s`           | Timeout duration            | 

### Run directly
```bash
sudo go run main.go -a 8.8.8.8 -id 2 -n 3
```

### Build executable
```bash
go build -o gopingid main.go
sudo ./gopingid -a 8.8.8.8 -id 2 -n 3
```

### Example Output
```bash
PING 8.8.8.8 (id=2):
8.8.8.8: icmp_seq=0 id=2 time=10ms
8.8.8.8: icmp_seq=1 id=2 time=18ms
8.8.8.8: icmp_seq=2 id=2 time=20ms
```

## License
This is free software under the terms of the [MIT License](https://github.com/yu-jimk/GoPingID/blob/main/LICENSE).
