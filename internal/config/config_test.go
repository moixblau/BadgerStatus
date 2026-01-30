package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Ensure env is clean for these vars
	envVars := []string{"VOLUMES", "SERIAL_PORT", "BAUD_RATE", "CRON", "HOST_PROC", "HOST_SYS"}
	old := make(map[string]string)
	for _, k := range envVars {
		old[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	// restore env after
	defer func() {
		for k, v := range old {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	cfg := LoadConfig()

	if !reflect.DeepEqual(cfg.Volumes, []string{"/"}) {
		t.Fatalf("expected volumes [/] got %v", cfg.Volumes)
	}
	if cfg.Port != "/dev/ttyACM0" {
		t.Fatalf("expected default port /dev/ttyACM0 got %s", cfg.Port)
	}
	if cfg.Rate != 115200 {
		t.Fatalf("expected default rate 115200 got %d", cfg.Rate)
	}
	if cfg.Cron != "@every 1m" {
		t.Fatalf("expected default cron '@every 1m' got %s", cfg.Cron)
	}
	if cfg.HostProc != "/proc" {
		t.Fatalf("expected default host_proc /proc got %s", cfg.HostProc)
	}
	if cfg.HostSys != "/sys" {
		t.Fatalf("expected default host_sys /sys got %s", cfg.HostSys)
	}
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	// backup and restore
	old := map[string]string{
		"VOLUMES":     os.Getenv("VOLUMES"),
		"SERIAL_PORT": os.Getenv("SERIAL_PORT"),
		"BAUD_RATE":   os.Getenv("BAUD_RATE"),
		"CRON":        os.Getenv("CRON"),
		"HOST_PROC":   os.Getenv("HOST_PROC"),
		"HOST_SYS":    os.Getenv("HOST_SYS"),
	}
	defer func() {
		for k, v := range old {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	os.Setenv("VOLUMES", "/data,/var")
	os.Setenv("SERIAL_PORT", "/dev/ttyUSB0")
	os.Setenv("BAUD_RATE", "9600")
	os.Setenv("CRON", "* * * * *")
	os.Setenv("HOST_PROC", "/my/proc")
	os.Setenv("HOST_SYS", "/my/sys")

	cfg := LoadConfig()

	if !reflect.DeepEqual(cfg.Volumes, []string{"/data", "/var"}) {
		t.Fatalf("expected volumes [/data /var] got %v", cfg.Volumes)
	}
	if cfg.Port != "/dev/ttyUSB0" {
		t.Fatalf("expected port /dev/ttyUSB0 got %s", cfg.Port)
	}
	if cfg.Rate != 9600 {
		t.Fatalf("expected rate 9600 got %d", cfg.Rate)
	}
	if cfg.Cron != "* * * * *" {
		t.Fatalf("expected cron '* * * * *' got %s", cfg.Cron)
	}
	if cfg.HostProc != "/my/proc" {
		t.Fatalf("expected host_proc /my/proc got %s", cfg.HostProc)
	}
	if cfg.HostSys != "/my/sys" {
		t.Fatalf("expected host_sys /my/sys got %s", cfg.HostSys)
	}
}
