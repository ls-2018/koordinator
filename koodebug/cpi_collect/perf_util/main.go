package main

import (
	"fmt"
	"github.com/hodgesds/perf-utils"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	containerCgroupFilePath := os.Args[1]
	//containerCgroupFilePath := "/sys/fs/cgroup/perf_event/kubepods.slice/kubepods-pod7bb341b7_92a7_4e41_919c_97b7e6fc68b7.slice/docker-48dcb0f76e8de43b8cd330b5e744aa92b54b47074712bd89cf382727b7999f54.scope"
	//containerCgroupFilePath := "/sys/fs/cgroup/perf_event/kubepods.slice/kubepods-pod7bb341b7_92a7_4e41_919c_97b7e6fc68b7.slice"
	//containerCgroupFilePath := "/sys/fs/cgroup/perf_event/kubepods.slice/"
	cgroupFile, err := os.OpenFile(containerCgroupFilePath, os.O_RDONLY, os.ModeDir)
	log.Println(containerCgroupFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	cpu := runtime.GOMAXPROCS(0) // 逻辑核
	ps := make([]perf.HardwareProfiler, cpu)
	for i := 0; i < cpu; i++ {
		cpiProfiler, err := perf.NewHardwareProfiler(int(cgroupFile.Fd()), i, perf.CpuCyclesProfiler|perf.CpuInstrProfiler, unix.PERF_FLAG_PID_CGROUP)
		if err != nil && !cpiProfiler.HasProfilers() {
			log.Fatalln(err)
		}
		go func() {
			if newErr := cpiProfiler.Start(); newErr != nil {
				log.Fatalln(err)
			}
		}()
		ps[i] = cpiProfiler

	}
	time.Sleep(time.Second * 10)
	for i := 0; i < cpu; i++ {
		profile := &perf.HardwareProfile{}
		if err := ps[i].Profile(profile); err != nil {
			log.Fatalln(err)
		}
		fmt.Println(i, *profile.CPUCycles, *profile.Instructions)
	}
	for i := 0; i < cpu; i++ {
		ps[i].Stop()
		ps[i].Close()
	}
	cgroupFile.Close()
}
