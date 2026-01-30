package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Volumes []string
	Port    string
	Rate    int
	Cron    string

	HostProc string
	HostSys  string

	Logger *slog.Logger
}

func LoadConfig() *Config {
	volumesStr := os.Getenv("VOLUMES")
	var volumes []string
	if volumesStr == "" {
		volumes = []string{"/"}
	} else {
		volumes = strings.Split(volumesStr, ",")
	}

	port := os.Getenv("SERIAL_PORT")
	if port == "" {
		port = "/dev/ttyACM0"
	}

	rateStr := os.Getenv("BAUD_RATE")
	if rateStr == "" {
		rateStr = "115200"
	}
	rate, err := strconv.Atoi(rateStr)
	if err != nil {
		log.Printf("Invalid BAUD_RATE: %v, using default 115200", err)
		rate = 115200
	}

	cronExpr := os.Getenv("CRON")
	if cronExpr == "" {
		cronExpr = "@every 1m"
	}

	hostProc := os.Getenv("HOST_PROC")
	if hostProc == "" {
		hostProc = "/proc"
	}

	hostSys := os.Getenv("HOST_SYS")
	if hostSys == "" {
		hostSys = "/sys"
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := &Config{
		Volumes:  volumes,
		Logger:   logger,
		Port:     port,
		Rate:     rate,
		Cron:     cronExpr,
		HostProc: hostProc,
		HostSys:  hostSys,
	}

	logger.Info("Configuration loaded",
		"volumes", cfg.Volumes,
		"port", cfg.Port,
		"rate", cfg.Rate,
		"cron", cfg.Cron,
		"host_proc", cfg.HostProc,
		"host_sys", cfg.HostSys)

	return cfg
}
