package main

import (
	"badgerStatus/internal/config"
	"badgerStatus/internal/metrics"
	"badgerStatus/internal/transport"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func sendMetrics(cfg *config.Config, bt *transport.BadgerTransport) {
	logger := cfg.Logger
	stats, err := metrics.GetStats(cfg.Volumes)
	if err != nil {
		logger.Error("Error getting stats", "volumes", cfg.Volumes, "error", err)
		return
	}

	msg := stats.String()

	_, err = bt.Write([]byte(msg))
	if err != nil {
		logger.Error("Write error", "error", err)
	}

	logger.Info("Metrics sent", "data", msg)
}

func main() {
	cfg := config.LoadConfig()
	logger := cfg.Logger

	// Configure gopsutil to use host's /proc and /sys
	os.Setenv("HOST_PROC", cfg.HostProc)
	os.Setenv("HOST_SYS", cfg.HostSys)

	bt, err := transport.NewBadgerTransport(cfg.Port, cfg.Rate)
	if err != nil {
		logger.Error("Failed to open serial port", "port", cfg.Port, "error", err)
		os.Exit(1)
	}
	defer bt.Close()

	logger.Info("Badger connected", "port", cfg.Port)

	// Wait to ensure Badger is ready
	time.Sleep(2 * time.Second)

	c := cron.New()
	_, err = c.AddFunc(cfg.Cron, func() {
		sendMetrics(cfg, bt)
	})
	if err != nil {
		logger.Error("Invalid CRON expression", "cron", cfg.Cron, "error", err)
		os.Exit(1)
	}
	c.Start()
	logger.Info("Cron started", "expression", cfg.Cron)
	select {} // Block forever
}
