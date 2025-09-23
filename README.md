# BBB Multicleaner

Automatically terminates long-running BigBlueButton meetings to prevent server overload.

## Installation

```bash
go build -o bbb-multicleaner
sudo mv bbb-multicleaner /usr/local/bin/
```

## Configuration

Edit these values in `main.go`:

```go
var MAX_DURATION = 4 * time.Hour     // Maximum meeting duration
var SLEEP_TIME = 15 * time.Minute    // Check interval
```

## Usage

### Manual Run

```bash
sudo bbb-multicleaner
```

### As Service

Create `/etc/systemd/system/bbb-multicleaner.service`:

```ini
[Unit]
Description=BBB Multicleaner
After=network.target

[Service]
ExecStart=/usr/local/bin/bbb-multicleaner
Restart=always
User=root

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable --now bbb-multicleaner
sudo systemctl status bbb-multicleaner
```

## How it works

1. Runs `bbb-conf --secret` to get BBB URL and secret
2. Calls BBB API to get active meetings
3. Ends meetings running longer than `MAX_DURATION`
4. Repeats every `SLEEP_TIME`

## Requirements

- BigBlueButton server
- Root access
- Go 1.19+ (for building)
