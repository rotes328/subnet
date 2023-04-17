# CLI Subnet Calculator
by Mike Rotella 2022

CLI tool for IP address subnetting written in Go.

To build:

```
$go build -o subnet main.go
$cp subnet /usr/local/bin
```

##Examples:

```
$subnet 10.10.10.10/29
╔══════════════════════════════════════════════════╗
║ For IP 10.10.10.10 and mask 255.255.255.248:     ║
╟──────────────────────────────────────────────────╢
║ Network Address:       10.10.10.8                ║
║ Broadcast Address:     10.10.10.15               ║
║ Range:                 10.10.10.9 - 10.10.10.14  ║
╚══════════════════════════════════════════════════╝
```


```
$subnet 192.168.17.49 255.255.255.192
╔══════════════════════════════════════════════════════╗
║ For IP 192.168.17.49 and mask 255.255.255.192:       ║
╟──────────────────────────────────────────────────────╢
║ Network Address:      192.168.17.0                   ║
║ Broadcast Address:    192.168.17.63                  ║
║ Range:                192.168.17.1 - 192.168.17.62   ║
╚══════════════════════════════════════════════════════╝
```
