package main

import (
	"fmt"
	"github.com/koordinator-sh/koordinator/apis/extension"
	"sort"
)

func main() {
	localCPUInfo, _ := GetLocalCPUInfo()

	nodeCPUInfo := &LocalCPUInfo{
		BasicInfo:      localCPUInfo.BasicInfo,
		ProcessorInfos: localCPUInfo.ProcessorInfos,
		TotalInfo:      localCPUInfo.TotalInfo,
	}
	cpus := make(map[int32]*extension.CPUInfo)
	cpuTopology := &extension.CPUTopology{}
	for _, cpu := range nodeCPUInfo.ProcessorInfos {
		info := extension.CPUInfo{
			ID:     cpu.CPUID,
			Core:   cpu.CoreID,
			Socket: cpu.SocketID,
			Node:   cpu.NodeID,
		}
		cpuTopology.Detail = append(cpuTopology.Detail, info)
		cpus[cpu.CPUID] = &info
	}
	sort.Slice(cpuTopology.Detail, func(i, j int) bool {
		return cpuTopology.Detail[i].ID < cpuTopology.Detail[j].ID
	})
	topo := NewCPUTopology(nodeCPUInfo) // ✅
	fmt.Println(topo)

	allCPUs := topo.CPUDetails.CPUs()
	reserved, _ := TakeByTopology(allCPUs, 3, topo) // ✈️ing

	fmt.Println(reserved)
}
