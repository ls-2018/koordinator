### https://help.aliyun.com/zh/alinux/user-guide/kernel-features/?spm=a2c4g.11186623.0.0.3227519eiptqLB
### hook

- batchresource 资源调整 从声明中获取调整后的资源
    - role
    - SetPodResources
    - SetContainerResources

- cpu_normalization 降低 limits.cpu QoSNone、QoSLS
    - role
    - AdjustPodCFSQuota 返回调整后的 limits.cpu * ratio TODO : 存不存在 这时 limits.cpu 比 requests.cpu 小
    - AdjustContainerCFSQuota 返回调整后的 limits.cpu * ratio TODO : 存不存在 这时 limits.cpu 比 requests.cpu 小

- cpuset
    - role
    - UnsetPodCPUQuota 设置了cpuset ，取消cfs_quota_us、cfs_period_us
    - SetContainerCPUSetAndUnsetCFS LSE, LSR pod 从声明中获取 CPUSET

    - beSharePools 是什么
    - sharePools 是什么


- gpu
    - InjectContainerGPUEnv 从声明中获取GPU设备信息，并注入 NVIDIA_VISIBLE_DEVICES

需要补一个 kubelet 的 none、static、

```
type Resources struct {
	CPUShares   *int64  // requests.Cpu().MilliValue()
	CFSQuota    *int64  // limits.Cpu().MilliValue()      有归一化处理
	CPUSet      *string
	MemoryLimit *int64  // limits.Memory().Value()
	CPUBvt *int64       // 声明中
}
```

```
逻辑差不多
cpu_normalization.go  adjustPodCFSQuota
batch_resources.go  SetPodCFSQuota

```

## typo

- kubepods.slice/kubepods-burstable.slice


- lscpu -e=CPU,NODE,SOCKET,CORE,CACHE,ONLINE
- lscpu -y
- /sys/devices/system/cpu/smt/active
- /sys/devices/system/cpu/intel_pstate/no_turbo
- /sys/fs/resctrl/info/L3/cbm_mask
- /proc/stat

socket 下 有 numa

？有什么区别
numa 32 c
socket 32 c
物理核 2 c

从sharedPoolCPUs中移除 (LSR && 没有在 annotations 中设置cpuSet) 但实际绑定了cpu的集合？

``` 逻辑核
sharedPoolCPUs = 节点所有资源的的拓扑 - 移除 (LSR|LSE && 没有在 annotations 中设置cpuSet) 但实际绑定的cpu - kubelet 预留的CPU - node 声明预留的CPU - node 系统使用的独占CPU

// coordinator定义的CPU共享池。该共享池主要供coordinator BE Pods或K8s Besteffort Pods使用。
BE ShardPool = sharedPoolCPUs - LSE使用的cpuSet 

// coordinator定义的CPU共享池。该共享池主要供coordinator LS Pods或K8s Burstable Pods使用。
LS ShardPool = sharedPoolCPUs - 所有使用的cpuSet 

sharePoolCPUCoresTotal = lscpu - (LSE、LSR)
sharePoolCPUCoresUsage = node_cpu_usage - (LSE、LSR、BE)
sharePoolUsageRatio    = sharePoolCPUCoresUsage / sharePoolCPUCoresTotal       
[0~0.9阈值,0.9阈值~阈值,阈值] = [nodeBurstIdle,nodeBurstCooling,nodeBurstOverload]
```

```
累计的CPU时钟滴答数指的是CPU在运行过程中所产生的时钟滴答数的总和。在计算机系统中，时钟滴答是指CPU时钟的一个周期，它用来衡量CPU的工作速度和时间。每个时钟周期都代表了CPU的一个基本操作步骤。

累计的CPU时钟滴答数可以用来衡量CPU的运行时间和负载情况。当CPU执行指令和处理任务时，时钟滴答数会不断累积。通过监测和记录累计的时钟滴答数，可以评估CPU的使用率、性能和效率。

在实际应用中，累计的CPU时钟滴答数可以用于分析系统的负载情况、优化程序的性能、检测CPU的瓶颈和判断系统资源的利用率。它是衡量系统和应用程序性能的重要指标之一。
```

lsblk -P -o NAME,TYPE,MAJ:MIN
lvs --noheadings --options lv_name,vg_name
vgs --noheadings --options vg_name,pv_count,pv_name
findmnt -P -o TARGET,SOURCE

#### cpu.stat

```
usage_usec：CPU使用的总时间，以微秒为单位。
user_usec：用户空间程序使用CPU的时间，以微秒为单位。
system_usec：内核空间程序使用CPU的时间，以微秒为单位。
core_sched.force_idle_usec：强制CPU空闲的时间，以微秒为单位。
nr_periods：在限制 CPU 使用的时间窗口内的周期数量。
nr_throttled：在限制 CPU 使用的时间窗口内被限制的周期数量。。
throttled_usec：在限制CPU使用的时间窗口内被限制的总时间，以微秒为单位。
nr_bursts：在限制CPU使用的时间窗口内的突发数量。
burst_usec：在限制CPU使用的时间窗口内的总突发时间，以微秒为单位。
```

- /sys/kernel/mm/kidled/use_hierarchy
- /sys/kernel/mm/kidled/scan_period_in_seconds
- /sys/kernel/mm/memcg_reaper/reap_background
- /sys/kernel/mm/memcg_reaper/min_free_kbytes
- /sys/kernel/mm/memcg_reaper/watermark_scale_factor
- memory.idle_page_stats
- /proc/meminfo

#### memUsageWithHotPageBytes
```
一个计算内存使用量的指标，它考虑了热页的字节数。

热页是指经常被访问的页，它们通常存储着常用的数据或代码。在计算内存使用量时，考虑热页的字节数可以更准确地反映实际的内存需求。

通常情况下，内存使用量可以通过统计所有页的字节数来计算。然而，这种方法没有考虑到热页的特殊性，可能会导致内存使用量的不准确估计。

memUsageWithHotPageBytes通过将热页的字节数纳入计算，提供了更精确的内存使用量指标。这意味着热页所占用的内存会被更好地考虑在内，从而更准确地反映出实际的内存需求。

总之，memUsageWithHotPageBytes是一个考虑了热页字节数的内存使用量指标，它提供了更准确的内存需求估计。

```

#### 龙蜥 memory.idle_page_stats

#### cgroup 版本
```
stat -fc %T /sys/fs/cgroup/
```



#### PSI、CPI 干扰检测

在Linux系统中，PSI（Pressure Stall Information）和CPI（Critical Pressure Information）是两个与系统性能和资源压力相关的指标。

1. PSI（Pressure Stall Information）：PSI是一种用于监控系统资源压力的机制，它通过跟踪不同资源（CPU、内存、磁盘、网络等）的压力情况，提供了更细粒度的资源利用信息。
   PSI将资源压力分为四个级别：空闲（idle）、干扰（pressure）、紧张（some）、严重（critical）。通过PSI指标，可以更好地了解系统资源是否受到压力影响，以及资源压力的严重程度。

2. CPI（Critical Pressure Information）：CPI是PSI的一个子集，它专门关注系统资源的严重压力情况。
   CPI提供了一种更高级别的资源压力报告，主要用于监控关键任务或关键服务是否受到资源压力的影响。
   CPI的目标是提供有关资源压力对关键任务的影响程度的信息，以便管理员可以及时采取措施来缓解资源压力。

通过使用PSI和CPI指标，系统管理员可以更好地了解系统资源的压力情况，并根据这些信息来优化系统性能、调整资源分配或处理资源瓶颈。这些指标对于监控和优化系统的性能是非常有用的。
- cpuacct/memory.pressure
- cpuacct/cpu.pressure
- cpuacct/io.pressure
https://zhuanlan.zhihu.com/p/656580184





```
	EvictByRealLimitPolicy   CPUEvictPolicy = "evictByRealLimit"   // 真正的限制 获取从历史获取到的 limit值
	EvictByAllocatablePolicy CPUEvictPolicy = "evictByAllocatable" // 可分配 从节点上获取 batch-cpu
```



#### 仅以下内核版本的Alibaba Cloud Linux镜像支持配置blk-iocost功能：
```
Alibaba Cloud Linux 2：4.19.81-17及以上内核版本。
Alibaba Cloud Linux 3：所有版本。
```


- /proc/sys/vm/min_free_kbytes  系统所保留空闲内存的最低限

https://www.zhihu.com/people/kang-kou-76

- cpu.cfs_burst_us
- cpu.cfs_quota_us

https://help.aliyun.com/zh/ack/ack-managed-and-ack-dedicated/user-guide/memory-qos-for-containers
https://developer.aliyun.com/article/976429?utm_content=m_1000349753
https://help.aliyun.com/zh/ack/ack-managed-and-ack-dedicated/getting-started/user-guide/?spm=a2c4g.11186623.0.0.6b115f1dV2GF3Q