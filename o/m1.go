package main

import (
	"context"
	"fmt"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {

	fmt.Println(eviction.ParseThresholdConfig([]string{}, map[string]string{
		"imagefs.available": "15%",
		"memory.available":  "100Mi",
		"nodefs.available":  "10%",
		"nodefs.inodesFree": "5%",
	}, nil, nil, nil))
	//lsCPUStr, _ := lsCPU("-e=CPU,NODE,SOCKET,CORE,CACHE,ONLINE")
	//processorInfos, _ := getProcessorInfos(lsCPUStr)
	//info := calculateCPUTotalInfo(processorInfos)
	//indent, _ := json.MarshalIndent(info, "  ", "   ")
	//fmt.Println(string(indent))
	//fmt.Println(string(indent))
}

type CPUTotalInfo struct {
	NumberCPUs  int32                     `json:"numberCPUs"`
	CoreToCPU   map[int32][]ProcessorInfo `json:"coreToCPU"`
	NodeToCPU   map[int32][]ProcessorInfo `json:"nodeToCPU"`
	SocketToCPU map[int32][]ProcessorInfo `json:"socketToCPU"`
	L3ToCPU     map[int32][]ProcessorInfo `json:"l3ToCPU"`
}

func calculateCPUTotalInfo(processorInfos []ProcessorInfo) *CPUTotalInfo {
	cpuMap := map[int32]struct{}{}
	coreMap := map[int32][]ProcessorInfo{}
	socketMap := map[int32][]ProcessorInfo{}
	nodeMap := map[int32][]ProcessorInfo{}
	l3Map := map[int32][]ProcessorInfo{}
	for i := range processorInfos {
		p := processorInfos[i]
		cpuMap[p.CPUID] = struct{}{}
		coreMap[p.CoreID] = append(coreMap[p.CoreID], p)
		socketMap[p.SocketID] = append(socketMap[p.SocketID], p)
		nodeMap[p.NodeID] = append(nodeMap[p.NodeID], p)
		l3Map[p.L3] = append(l3Map[p.L3], p)
	}
	return &CPUTotalInfo{
		NumberCPUs:  int32(len(cpuMap)),
		CoreToCPU:   coreMap,
		SocketToCPU: socketMap,
		NodeToCPU:   nodeMap,
		L3ToCPU:     l3Map,
	}
}
func lsCPU(option string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	executable, err := exec.LookPath("lscpu")
	if err != nil {
		return "", fmt.Errorf("failed to lookup lscpu path, err: %w", err)
	}
	output, err := exec.CommandContext(ctx, executable, option).Output()
	if err != nil {
		return "", fmt.Errorf("failed to exec command %s, err: %v", executable, err)
	}
	return string(output), nil
}

type ProcessorInfo struct {
	// logic CPU/ processor ID
	CPUID int32 `json:"cpu"`
	// physical CPU core ID
	CoreID int32 `json:"core"`
	// cpu socket ID
	SocketID int32 `json:"socket"`
	// numa node ID
	NodeID int32 `json:"node"`
	// L1 L2 cache ID
	L1dl1il2 string `json:"l1dl1il2"`
	// L3 cache ID
	L3 int32 `json:"l3"`
	// online
	Online string `json:"online"`
}

func getProcessorInfos(lsCPUStr string) ([]ProcessorInfo, error) {
	if len(lsCPUStr) <= 0 {
		return nil, fmt.Errorf("lscpu output is empty")
	}

	var processorInfos []ProcessorInfo
	for _, line := range strings.Split(lsCPUStr, "\n") {
		items := strings.Fields(line)
		if len(items) < 6 {
			continue
		}
		cpu, err := strconv.ParseInt(items[0], 10, 32)
		if err != nil {
			continue
		}
		node, _ := strconv.ParseInt(items[1], 10, 32)
		socket, err := strconv.ParseInt(items[2], 10, 32)
		if err != nil {
			continue
		}
		core, err := strconv.ParseInt(items[3], 10, 32)
		if err != nil {
			continue
		}
		l1l2, l3, err := GetCacheInfo(items[4])
		if err != nil {
			continue
		}
		online := strings.TrimSpace(items[5])
		info := ProcessorInfo{
			CPUID:    int32(cpu),
			CoreID:   int32(core),
			SocketID: int32(socket),
			NodeID:   int32(node),
			L1dl1il2: l1l2,
			L3:       l3,
			Online:   online,
		}
		processorInfos = append(processorInfos, info)
	}
	if len(processorInfos) <= 0 {
		return nil, fmt.Errorf("no valid processor info")
	}

	// sorted by cpu topology
	// NOTE: in some cases, max(cpuId[...]) can be not equal to len(processors)
	sort.Slice(processorInfos, func(i, j int) bool {
		a, b := processorInfos[i], processorInfos[j]
		if a.NodeID != b.NodeID {
			return a.NodeID < b.NodeID
		}
		if a.SocketID != b.SocketID {
			return a.SocketID < b.SocketID
		}
		if a.CoreID != b.CoreID {
			return a.CoreID < b.CoreID
		}
		return a.CPUID < b.CPUID
	})

	return processorInfos, nil
}
func GetCacheInfo(str string) (string, int32, error) {
	// e.g.
	// $ `lscpu -e=CPU,NODE,SOCKET,CORE,CACHE,ONLINE`
	// CPU NODE SOCKET CORE L1d:L1i:L2:L3 ONLINE
	//  0    0      0    0 0:0:0:0          yes
	//  1    0      0    0 0:0:0:0          yes
	//  2    0      0    1 1:1:1:0          yes
	//  3    0      0    1 1:1:1:0          yes
	infos := strings.Split(strings.TrimSpace(str), ":")
	// assert l1, l2 are private cache, so they have the same id with the core
	// L3 cache maybe not available, when the host is qemu-kvm. detail: https://bugzilla.redhat.com/show_bug.cgi?id=1434537
	if len(infos) < 3 {
		return "", 0, fmt.Errorf("invalid cache info %s", str)
	}
	l1l2 := infos[0]
	if len(infos) == 3 {
		return l1l2, 0, nil
	}
	l3, err := strconv.ParseInt(infos[3], 10, 32)
	if err != nil {
		return "", 0, err
	}
	return l1l2, int32(l3), nil
}
