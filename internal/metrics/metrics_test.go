package metrics

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func TestStats_String(t *testing.T) {
	s := Stats{CPUUsage: 10, RAMUsed: 100, RAMTotal: 200, DiskUsed: 300, DiskTotal: 400}
	sstr := s.String()
	if !strings.Contains(sstr, "cpu=10") {
		t.Fatalf("unexpected string: %s", sstr)
	}
	if !strings.Contains(sstr, "ram_used=100") {
		t.Fatalf("unexpected string: %s", sstr)
	}
}

func TestGetStats_Success(t *testing.T) {
	// save originals and restore
	origCPU := cpuPercent
	origMem := memVirtualMemory
	origDisk := diskUsage
	defer func() { cpuPercent = origCPU; memVirtualMemory = origMem; diskUsage = origDisk }()

	// inject fakes
	cpuPercent = func(interval time.Duration, percpu bool) ([]float64, error) {
		return []float64{12.0}, nil
	}
	memVirtualMemory = func() (*mem.VirtualMemoryStat, error) {
		return &mem.VirtualMemoryStat{Used: 50 * 1024 * 1024, Total: 100 * 1024 * 1024}, nil
	}
	diskUsage = func(path string) (*disk.UsageStat, error) {
		return &disk.UsageStat{Used: 20 * 1024 * 1024, Total: 100 * 1024 * 1024}, nil
	}

	stats, err := GetStats([]string{"/fake"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.CPUUsage != 12 {
		t.Fatalf("expected CPU 12 got %d", stats.CPUUsage)
	}
	if stats.RAMUsed != 50 || stats.RAMTotal != 100 {
		t.Fatalf("unexpected RAM values: used=%d total=%d", stats.RAMUsed, stats.RAMTotal)
	}
	if stats.DiskUsed != 20 || stats.DiskTotal != 100 {
		t.Fatalf("unexpected Disk values: used=%d total=%d", stats.DiskUsed, stats.DiskTotal)
	}
}

func TestGetStats_DiskError(t *testing.T) {
	origCPU := cpuPercent
	origMem := memVirtualMemory
	origDisk := diskUsage
	defer func() { cpuPercent = origCPU; memVirtualMemory = origMem; diskUsage = origDisk }()

	cpuPercent = func(interval time.Duration, percpu bool) ([]float64, error) {
		return []float64{1.0}, nil
	}
	memVirtualMemory = func() (*mem.VirtualMemoryStat, error) {
		return &mem.VirtualMemoryStat{Used: 1 * 1024 * 1024, Total: 2 * 1024 * 1024}, nil
	}
	diskUsage = func(path string) (*disk.UsageStat, error) {
		return nil, errors.New("disk error")
	}

	_, err := GetStats([]string{"/fake"})
	if err == nil {
		t.Fatalf("expected error from diskUsage, got nil")
	}
}
