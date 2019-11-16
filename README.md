# LED-Simple-Monitoring

## Installation
Requires a supported release of Go. 
built using `go 1.12.7`

`go get -u github.com/alonchis/LED-Simple-Monitoring`

## Usage
Edit led-simple-monitoring.go **pinsIndex** to match your setup
**sites** to define which sites to monitor, and **period** to set the interval in which to check the sites

Once those are defined, run
```bash
sudo go run led-simple-monitoring.go > /var/log/led.log 2>&1 & 
```
logs will be written to /var/log/led.log and run in the background. 


### Running as system service
All commands below run as user `pi
`
build executable and add to system bin
```bash
$> sudo go build led-simple-monitorig.go
$> sudo cp ./led-simple-monitoring /usr/bin
    
```

copy `ledmonitoring.service` to /etc/systemd/system 
```bash
$> sudo cp ledmonitorin.service /etc/systemd/system
$> sudo systemctl enable ledmonitoring.service
$> sudo systemctl start ledmonitoring.service
```
after this, led-monitoring should run automatically even after restart

### stopping services
if need to stop service run `sudo systemctl status`

if run as background process, run `ps aux | grep -i led` 
to find PID to kill via `sudo kill -0 PID`