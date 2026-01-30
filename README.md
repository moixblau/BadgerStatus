# BadgerStatus

A small system status monitor for the Badger 2040 e-ink badge.

This project collects simple system metrics (CPU, RAM, disk) and sends them to a connected Badger 2040 over serial so the device can display them. It is designed to be run on a host and the badge connected by USB.

## Recommended: Run with Docker Compose

The easiest way to run this on a NAS (or any Linux host with Docker) is with Docker Compose.

1. Copy `docker-compose.yml` to your host.
2. Make sure the Badger 2040 is connected by USB (commonly appears as `/dev/ttyACM0`).
3. Start the service:

## Environment variables

These variables are read by the application at startup. The defaults match common setups but can be overridden via environment or in a Docker Compose file.

| Variable      | Description                                                                                  |        Default |
|---------------|----------------------------------------------------------------------------------------------|---------------:|
| `VOLUMES`     | Comma-separated list of mount points to monitor (disk usage).                                |            `/` |
| `SERIAL_PORT` | Serial device path for the Badger 2040.                                                      | `/dev/ttyACM0` |
| `BAUD_RATE`   | Baud rate used for serial communication.                                                     |       `115200` |
| `CRON`        | Cron expression or schedule for when metrics are sent (e.g. `@every 1m` or `*/5 * * * *`).   |    `* * * * *` |
| `HOST_PROC`   | Path to host `/proc` (mounted into container).                                               |        `/proc` |
| `HOST_SYS`    | Path to host `/sys` (mounted into container).                                                |         `/sys` |

Note: The code reads `VOLUMES` (plural) — please use that exact name in your environment or compose file.

## Tests

This repository contains unit tests for core packages. Run them with:

```bash
go test ./...
```

Tests include:
- `internal/config` — checks default and env-based configuration.
- `internal/metrics` — unit tests mock system calls to validate metric calculations.
- `internal/transport` — unit tests mock the serial port so tests run without hardware.

## Development notes

- Metrics collection uses `github.com/shirou/gopsutil` and is implemented in `internal/metrics`.
- The serial transport is in `internal/transport` and uses `go.bug.st/serial` — during tests the package-level `serialOpen` seam is replaced by a fake implementation to avoid needing actual hardware.
- The `forBadger2040/monitor.py` script is a simple client for the badge. It reads metric lines from stdin (format: `cpu=... ram_used=... ram_total=... disk_used=... disk_total=...`) and paints them on the Badger. You can run it on the device or emulate by piping lines to it.