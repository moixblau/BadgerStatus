package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// package-level function variables so tests can replace these with fakes
var (
	cpuPercent       = cpu.Percent
	memVirtualMemory = mem.VirtualMemory
	diskUsage        = disk.Usage
)

type Stats struct {
	CPUUsage  int
	RAMUsed   uint64
	RAMTotal  uint64
	DiskUsed  uint64
	DiskTotal uint64
}

func (s Stats) String() string {
	return fmt.Sprintf(
		"cpu=%d ram_used=%d ram_total=%d disk_used=%d disk_total=%d\n",
		s.CPUUsage, s.RAMUsed, s.RAMTotal, s.DiskUsed, s.DiskTotal,
	)
}

func GetStats(volumes []string) (Stats, error) {
	// CPU
	cpuPerc, err := cpuPercent(0, false)
	if err != nil {
		return Stats{}, err
	}
	cpuPct := int(cpuPerc[0])

	// RAM
	vm, err := memVirtualMemory()
	if err != nil {
		return Stats{}, err
	}
	ramUsed := vm.Used / (1024 * 1024)
	ramTotal := vm.Total / (1024 * 1024)

	// DISK (total of all volumes)
	var diskUsed, diskTotal uint64
	for _, volume := range volumes {
		du, err := diskUsage(volume)
		if err != nil {
			return Stats{}, err
		}
		diskUsed += du.Used / (1024 * 1024)
		diskTotal += du.Total / (1024 * 1024)
	}

	return Stats{
		CPUUsage:  cpuPct,
		RAMUsed:   ramUsed,
		RAMTotal:  ramTotal,
		DiskUsed:  diskUsed,
		DiskTotal: diskTotal,
	}, nil
}
