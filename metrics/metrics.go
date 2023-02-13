package metrics

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/c9s/goprocinfo/linux"
	"github.com/nano-ci-cd/util"
	"github.com/robfig/cron/v3"
)

type MetricsConfig struct {
}

type MemoryEntity struct {
	Free     int
	Used     int
	SwapFree int
	SwapUsed int
}

type StorageEntity struct {
	Total int
	Free  int
}

func Start() {
	c := cron.New()

	c.AddFunc("@every 1s", func() {
		stats, err := getStorageStats()

		if err != nil {
			println(err.Error())
		}

		println("Storage stats:" + fmt.Sprintf("Total: %v, Free: %v", stats.Total, stats.Free))
	})

	go c.Start()
}

// todo for later below

func getStorageStats() (StorageEntity, error) {
	total, err := parseStorageCommand(storageTotal)

	if err != nil {
		return StorageEntity{}, err
	}

	free, err := parseStorageCommand(storageFree)

	if err != nil {
		return StorageEntity{}, err
	}

	return StorageEntity{
		Total: total,
		Free:  free,
	}, nil
}

func parseStorageCommand(commandOutput string) (int, error) {
	out, err := util.ExecuteCommandWithOutput(commandOutput)

	if err != nil {
		return 0, err
	}

	a, err := strconv.Atoi(out)

	if err != nil {
		return 0, err
	}

	return a, nil
}

func getMemoryStats() (MemoryEntity, error) {
	memTotal, memUsed, err := parseMemCommand(memCommand)

	if err != nil {
		return MemoryEntity{}, err
	}

	swapTotal, swapUsed, err := parseMemCommand(swapCommand)

	if err != nil {
		return MemoryEntity{}, err
	}

	return MemoryEntity{
		Free:     memTotal,
		Used:     memUsed,
		SwapFree: swapTotal,
		SwapUsed: swapUsed,
	}, nil
}

func parseMemCommand(commandOutput string) (int, int, error) {
	out, err := util.ExecuteCommandWithOutput(commandOutput)

	if err != nil {
		return 0, 0, err
	}

	splitted := strings.Split(out, " ")

	if len(splitted) != 2 {
		return 0, 0, errors.New("could not parse memory command output")
	}

	a, err := strconv.Atoi(splitted[0])

	if err != nil {
		return 0, 0, err
	}

	b, err := strconv.Atoi(splitted[1])

	if err != nil {
		return 0, 0, err
	}

	return a, b, nil
}

func getCpuInfo() error {
	a, err := linux.ReadLoadAvg("/proc/loadavg")

	if err != nil {
		return err
	}

	println(a.Last1Min)
	diskStat, err := linux.ReadDisk("/proc/diskstats")

	if err != nil {
		return err
	}

	println(fmt.Sprintf("Disk usage: %v/%v", diskStat.Used, diskStat.All))

	return nil
}

func parseCpuCommand(commandOutput string) (float64, error) {
	out, err := util.ExecuteCommandWithOutput(commandOutput)

	if err != nil {
		return 0, err
	}

	a, err := strconv.ParseFloat(out, 64)

	if err != nil {
		return 0, err
	}

	return a, nil
}

const (
	memCommand    = "free | grep Mem | awk '{print $2,$3}'"
	swapCommand   = "free | grep Swap | awk '{print $2,$3}'"
	storageTotal  = "df | awk '{print $2}' | tail -n +2 | paste -sd+ | bc"
	storageFree   = "df | awk '{print $3}' | tail -n +2 | paste -sd+ | bc"
	cpuPercentage = "top -b -n 1 | awk '{print $9}' | tail -n +8 | paste -sd+ | sed s/,/./g  | bc"
)
