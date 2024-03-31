cgroup2 
- https://zhuanlan.zhihu.com/p/637248596
- https://zorrozou.github.io/
- https://zhuanlan.zhihu.com/p/654673360


grpc
- https://www.cnblogs.com/ssezhangpeng/p/12402189.html

```
kubectl proxy --port=8001
curl http://localhost:8001/api/v1/nodes/kind-worker/proxy/configz
curl http://localhost:8001/api/v1/nodes/kind-worker/proxy/pods
```


// Node(Batch).Alloc[usage] := Node.Total - Node.Reserved - System.Used - sum(Pod(Prod/Mid).Used)
// System.Used = max(Node.Used - Pod(All).Used, Node.Anno.Reserved, Node.Kubelet.Reserved)
// Node(Batch).Alloc[request] := Node.Total - Node.Reserved - System.Reserved - sum(Pod(Prod/Mid).Request)
// System.Reserved = max(Node.Anno.Reserved, Node.Kubelet.Reserved)
// Node(Batch).Alloc[maxUsageRequest] := Node.Total - Node.Reserved - System.Used - sum(max(Pod(Prod/Mid).Request, Pod(Prod/Mid).Used))



- CPU 绑定策略
  - kubelet
    - full-pcpus-only 确保容器分配的 CPU 物理核独占
    - distribute-cpus-across-numa 按照物理核维度均匀的分配CPU
  - koord
    - CPUBindPolicyDefault
    - CPUBindPolicyFullPCPUs 与 full-pcpus-only 类似(尽可能确保容器分配的 CPU 物理核独占)
    - CPUBindPolicySpreadByPCPUs 与 distribute-cpus-across-numa 一样
    - CPUBindPolicyConstrainedBurst 与 distribute-cpus-across-numa 类似,但会从对应POOL中获取
- CPU 独占策略
  - CPUExclusivePolicyDefault
  - CPUExclusivePolicyPCPULevel     尽量避开已经被同一个独占策略申请的物理核。它是对 CPUBindPolicySpreadByPCPUs 策略的补充。
  - CPUExclusivePolicyNUMANodeLevel 在分配逻辑 CPU 时，尽量避免 NUMA 节点已经被相同的独占策略申请。如果没有满足策略的 NUMA 节点，则降级为 CPUExclusivePolicyPCPULevel 策略
- NUMA 拓扑对齐策略  (容器级、不同资源尽量对齐到一个numa)  
  - kubelet
    - None
    - BestEffort                    优先选择拓扑对齐的 NUMA Node，如果没有，则继续为 Pods 分配资源
    - Restricted                    每个 Pod 在 每个NUMA 节点上请求的资源是拓扑对齐的，如果不是，koord-scheduler 会在调度时跳过该节点
    - SingleNUMANode                表示一个 Pod 请求的所有资源都必须在同一个 NUMA 节点上，如果不是，koord-scheduler 调度时会跳过该节点。
- NUMA 分配策略
  - kubelet   以kubelet cpu绑定策略优先 
    - MostAllocated
    - LeastAllocated
    - DistributeEvenly
  - bin-packing
    - 分配最空闲的 NUMA 节点

忽略 kubelet 预留

- CPU Shared Pool     动态设置 cgroup         包含节点中所有未分配的CPU;不包括由 K8s Guaranteed、LSE 和 LSR Pod 分配的 CPU
  - K8s Burstable 
  - LS  
  - K8s Guaranteed 分数 CPU requests
  - 当LSE、LSR pod 创建 销毁时, 调整CPU Shared Pool  

- BE CPU Shared Pool   动态设置 cgroup      包含节点中除 K8s Guaranteed 和 Koordinator LSE Pod 分配的之外的所有 CPU
  - K8s BestEffort 
  - BE 

- 静态独占 (从CPU Shared Pool 申请、归还)
  - LSE
    - 完全独占
    - K8s Guaranteed 整数 CPU requests
  - LSR
    - 可以共享到BE POOL
    - K8s Guaranteed 整数 CPU requests

 
- LS
  - 直接标注 Koordinator QoS && K8s Guaranteed Pod && cpu是1000的整数倍
  - 没有标注 Koordinator QoS && K8s Guaranteed Pod && kubelet 的 CPU 管理器策略为 none 策略 && cpu是1000的整数倍
  - 绑定 LSE、LSR独占之外的所有CPU
- LSR
  - 直接标注 Koordinator QoS && K8s Guaranteed Pod && cpu是1000的整数倍
  - 没有标注 Koordinator QoS && K8s Guaranteed Pod && kubelet 的 CPU 管理器策略为 static 策略 && cpu是1000的整数倍
  - 可以与BE共享
  - 不会从kubelet 预留的cpu中分配
- LSE
  - 逻辑CPU是完全独占的，不得共享
  - 不会从kubelet 预留的cpu中分配
- BE
  - 绑定 LSE、LSR独占之外的所有CPU



- 1592
- 1600
- 1603
- 1703
- 1768
- 1774
- 1788