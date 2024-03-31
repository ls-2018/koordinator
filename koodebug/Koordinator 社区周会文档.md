# 日程

|时间|（每双周）周二晚上 19:30 - 20:30 北京时间|
|:----|:----|
|钉钉会议|入会链接：[https://meeting.dingtalk.com/j/XG2v8SPPA5u](https://meeting.dingtalk.com/j/XG2v8SPPA5u)|
|钉群提醒/入会|搜索钉钉群号 33383887|
|往期视频|[https://space.bilibili.com/3493089995917452/video](https://space.bilibili.com/3493089995917452/video)|

# 议题&记录


---
## 2023/12/20 #21

SKIP

## 2023/12/05 #20

主持人：@jasonliu747

参会人：

议题：

1. Koordinator YARN Copilot使用： 
2. Network QoS的相关近况
3. Mid tier proposal讨论 [https://github.com/koordinator-sh/koordinator/pull/1762](https://github.com/koordinator-sh/koordinator/pull/1762)
4. PodMigrationJob support eviction conditions proposal 讨论 [https://github.com/koordinator-sh/koordinator/pull/1767](https://github.com/koordinator-sh/koordinator/pull/1767)
## 
## 2023/11/07 #19

主持人：@zwzhang0107

参会人：

议题：

1. Mid tier 技术演进讨论：[https://shimo.im/docs/5xkGoegGB8fvM9kX](https://shimo.im/docs/5xkGoegGB8fvM9kX)（[https://shimo.im/file-invite/nYVgqwEv9f9M5a66JadyoR6pQXpb6/](https://shimo.im/file-invite/nYVgqwEv9f9M5a66JadyoR6pQXpb6/)）
2. Network QoS方案介绍 [https://shimo.im/docs/XKq42Lv6nOugglAN](https://shimo.im/docs/XKq42Lv6nOugglAN) (https://shimo.im/file-invite/AVjodR3Uv3qc4vGAca78aqJwyY936/ l1b0k)
3. TODO
## 2023/10/24 #18

主持人:  @FillZpp

参会人:

议题:

1. [Hadoop YARN on K8s](https://koordinator.sh/zh-Hans/docs/next/best-practices/colocation-of-hadoop-yarn/) 版本能力预告
2. Network QoS（Terway）共建合作
3. 1.4 版本发布时间
4. Open(RDT plugin usage)
## 2023/10/10

SKIP

## 2023/09/12 #17

主持人:  @eahydra

参会人:

议题:

1. CPU绑核策略支持Required语义gith
2. NUMA-Aware Scheduling 进展
3. Reservation 相关的变化：框架层面支持Reservation维度的打分归一化；另外 Reservation 调度后，资源汇总方式发生了变化，几点考虑：一个是方便后续做Reservation维度的规格变配，另一个是预留GPU资源后，可以支持用Reservation资源时使用GPU Share能力。
## 
## 2023/08/23 #16

主持人:

参会人:

议题:

1. 社区 v1.3 版本 release note 介绍: [https://koordinator.sh/blog/release-v1.3.0/](https://koordinator.sh/blog/release-v1.3.0/)
2. v1.4 版本规划介绍: [https://github.com/koordinator-sh/koordinator/milestone/12](https://github.com/koordinator-sh/koordinator/milestone/12)
3. 关于在ElasticQuota中添加koord-root-quota的开发方案讨论。Issue连接：[https://github.com/koordinator-sh/koordinator/issues/1519](https://github.com/koordinator-sh/koordinator/issues/1519)

## 2023/08/08 #15

主持人：zwzhang0107

参会人：张左玮、逐灵、申信、乔普、吕风、酒祝、佳风、鲍文杰、几明、刘明、王兴奇、bowen、kangzhang、CH

议题:

1. 社区 v1.3 版本 issue 进展对齐
2. Kubecon China峰会Koordinator议题介绍
3. 开源之夏议题proposal介绍
    1. 重调度器支持仲裁机制提升驱逐稳定性（[proposal: Eviction Arbitration Mechanism in Descheduler by baowj-678 · Pull Request #1454 · koordinator-sh/koordinator (github.com)](https://github.com/koordinator-sh/koordinator/pull/1454)）
    2. 支持冷内存的采集和上报 （[https://github.com/koordinator-sh/koordinator/pull/1514](https://github.com/koordinator-sh/koordinator/pull/1514)）
4. CPU归一化和资源超卖介绍 @佳风
5. 怎么衡量应用被干扰

---
## 2023/07/11 #14

主持人：

参会人：

议题:

1. 社区 v1.3 版本 issue 进展对齐
2. 提议：是否可以在 ElasticQuota 方案中添加不可抢占的特性 （by 胜同学）
    1. 背景：
        1. Google Borg 中对于 Production 负载有个特性定义，该等级的负宋博文、王兴奇、张康、周岳骞、CH、j4ckstraw、ryan、zjk载不可以被抢占(Although a preempted task will often be rescheduled elsewhere in the cell, preemption cascades could occur if a high-priority task bumped out a slightly lower-priority one, which bumped out another slightly-lower priority task, and so on. To eliminate most of this, **we disallow tasks in the production priority band to preempt one another.** )
        2. 在 [Multi Hierarchy Elastic Quota Management | Koordinator](https://koordinator.sh/docs/designs/multi-hierarchy-elastic-quota-management) 方案中，如果某个 Quota 下的负载使用了超过 Min 的资源，无论其 Priority/QoS 等级如何，该负载可能会被抢占；
    2. 提议：对于重要负载，是否可以支持设置 Pod 为不可抢占（比如设置 koordinator.sh/preemptible: false），那么设置了该特性的 Pod，限定其资源使用之和就不能超过 Quota 的 Min，这样就不会被抢占；（或者可以通过配置开启限定 koord-prod Priority的负载不能使用超过 Min 的资源）
    3. 参考：Uber Peloton，[https://peloton.readthedocs.io/en/latest/user-guide/#resource-pools](https://peloton.readthedocs.io/en/latest/user-guide/#resource-pools)
3. 期望了解下资源画像的开发排期
4. Performance collector perf plugin改进（CPI use cycles instead of ref-cycles?）

---
## 2023/06/27 #13

主持人：酒祝

参会人：

议题:

5. mid tier 议题分享，以及未来对于在线超卖的一些构想
    1. [https://github.com/koordinator-sh/koordinator/pull/1385](https://github.com/koordinator-sh/koordinator/pull/1385)
6. scheduler 框架升级 1.24 进展，以及兼容处理机制
7. 应用资源画像有没有计划？前置解决热点节点问题

---
## 2023/06/13 #12

主持人：Jason

参会人：吕风、逐灵、酒祝、张佐玮、乔普、申信

议题:

8. Koordinator on NRI 方案探讨 @张康 (Intel)
    1. [https://github.com/koordinator-sh/koordinator/pull/1366](https://github.com/koordinator-sh/koordinator/pull/1366)
9. 

---
## 2023/05/30 #11

主持人：吕风

参会人：申信、吕风、逐灵、酒祝、张左玮、乔普

议题：

1. Github Project 需求功能同步
2. 评审通过的Proposal: NUMA Topology Scheduling
3. koordinator-yarn v1.0 同步
    1. 提问：关于koordinator-yarn的技术方案
    2. 如何提供磁盘和网络的 QoS 保障，避免离线负载的干扰
    3. 如何管理 Spark  的 Shuffle 数据
4. 提问：是否可以调整 Koord PriorityClass 的默认值，比如 koord-prod  默认值为 0
    1. 背景：K8s 集群部署 Koordinator 过程中，存量 Pod 的 Priority 为 0，低于 Koord 定义的 Priority（koord-prod/koord-middle），导致存量 Pod 有被抢占的风险；
5. koord-scheduler近期的变化
    1. Reservation 支持 AllocatePolicy，可以避免Pod复用多个Reservation资源的问题。
    2. 新增 Reservation Affinity 协议，支持在 Pod 上声明要求使用哪些 Reservation。当集群内没有符合 Reservation Affinity 的Reservation 时，Pod会调度失败；如果存在匹配的Reservation，则会从这些Reservation中分配。该特性同时会检查 Reservation 声明的资源Owner是否与Pod匹配。
    3. 引入 Reservation 后导致 k8s 原生的scheduler plugins 打分异常问题已经得到修复。
    4. 性能优化：调度器PreBind阶段各个插件的Patch Pod行为从多次收敛为1次。
6. koord-manager 近期的变化
    1. ClusterColocationProfile 支持新的灰度机制，支持设置百分比，把一部分Pod切换为BE QosClass
7. koordlet 近期的变化
    1. 指标存储切换为tsdb

---
## 2023/04/18 #10

主持人：申信

参会人：申信、吕风、逐灵、酒祝、张左玮、乔普、高瑞鸿、孙骁、CH、宋涛、吕波

议题：

8. Github Project 需求功能同步
9. Reservation 增强进展介绍
10. 


---
## 2023/04/04 #9

主持人：醒宝

参会人：醒宝、酒祝、吕风、Jun、CH、死磕、张左玮、刘明、宋泽辉、高瑞鸿、乔普、申信、fjding、逐灵、叶欢欢、宋涛

议题：

1. reservationg功能\pr分享
2. AMD LLC/MBA 隔离特性介绍
3. 1.2 版本发布信息同步
4. 探讨 NRI v2 作为 runtime hook 触发源的可行性

# 议题&记录


---
## 2023/03/21 #8

主持人：Jason

参会人：Jason、酒祝、逐灵、乔普、左峰、怀仁、高瑞鸿、江博、张左玮、岳溟、within、申信、张凯、杨峻峰、lrving、CH、胡伟、刘明

议题：

1. GPU拓扑感知调度 [https://github.com/koordinator-sh/koordinator/pull/1115](https://github.com/koordinator-sh/koordinator/pull/1115)
2. 

---
## 2023/03/07 #7

主持人：酒祝

参会人：TBD

议题：

3. Koordinator v1.2 版本计划
    1. reservation 能力补全
    2. node resource reservation 进展
4. descheduler 支持 Kruise PUB 防护（[PR](https://github.com/koordinator-sh/koordinator/pull/1064)），Koordinator & OpenKruise 社区联合讨论
5. 干扰检测进展同步

---
## 2023/02/21 #6

主持人：逐灵

参会人：逐灵、酒祝、张佐玮、刘明、吕风、宋涛、Jason、CH、程兴源、刘岩、乔普、申信、怀仁、死磕、高瑞鸿、Te、JF、岳溟、左峰、66、王杰诚

议题：

1. Koordinator v1.2 版本研发进展同步
2. 干扰率相关的事情除了目前暴露的cpi和psi，还可以考虑哪些纬度
3. network io隔离的部分，已经有bandwidth类似的containernetwork plugin在做，是否还有必要在koordinator中实现类似的功能
4. koordetector干扰检测计划
    1. [https://github.com/koordinator-sh/koordetector/pull/3/files](https://github.com/koordinator-sh/koordetector/pull/3/files)
    2. [https://github.com/koordinator-sh/koordetector/pull/4](https://github.com/koordinator-sh/koordetector/pull/4)
5. koordlet指标采集框架化改造
    1. [https://github.com/koordinator-sh/koordinator/pull/1019](https://github.com/koordinator-sh/koordinator/pull/1019)
6. 社区考虑怎么支持flink的任务调度
## 2023/02/07 #5

主持人：吕风

参会人：申信、逐灵、吕风、酒祝、Jason、高瑞鸿、张佐玮、乔普、吴昆、宋涛、刘明、左峰、海馨、CH、辰、立衡、魏博锴、miea、Tendrun、江博、张薇薇、死磕、J4ckstraw、谢建超

议题：

1. 最近 Koordinator 的进展和 v1.2 版本规划
2. 讨论runtime-proxy两种模式的优劣势  @刘岩
3. ![image-20240229150847776](https://raw.githubusercontent.com/FlyFishking/picgo/main/liushuo/202402291508808.png)


3. bvt + noise clean on centos @扶风（龙蜥团队）
4. cgroups v2 如何支持流量打标

---
## 2023/01/10 (Cancel)


---
## 2023/01/10 #4

主持人：申信

参会人：申信、逐灵、吕风、酒祝、刘明、高瑞鸿、张佐玮、张同学、乔普、立衡、刘岩、车漾、吴昆、王泽360、张凯、怀仁、lrving、小月半子、Jason、zichen、宋涛、谢建超、博易、J4ckstraw、CH、Tendrun

议题：

5. Koordinator 重调度功能讨论
    1. 迁移防御策略（[#933](https://github.com/koordinator-sh/koordinator/pull/933)）
6. Koordinator 静态资源预留 proposal 介绍 & 讨论（[#922](https://github.com/koordinator-sh/koordinator/pull/922)）
7. Memcg metrics采集的方式 (这是一类问题，龙蜥内核特性的metrics使用什么方式采集比较好)

---
## 2022/12/27 #3

主持人：酒祝

参会人：酒祝、逐灵、佑祎、吕风、申信、xzy、j4ckstraw、王孝诚、刘岩、高瑞鸿、乔普、吴昆、宋涛、Herbertduan、睿元、明昼、刘明

议题：

1. Koordinator v1.1 release 版本介绍
    1. 负载感知调度&重调度
    2. 指标采集、干扰检测
2. 干扰检测中使用eBPF收集CPU调度延迟 proposal 介绍&讨论 （[#859](https://github.com/koordinator-sh/koordinator/pull/859)、[#860](https://github.com/koordinator-sh/koordinator/pull/860)）

---
## 2022/12/13 #2

主持人：酒祝

参会人：酒祝、逐灵、佑祎、吕风、CH、姚翔、Herbert、魏博锴、乔普、Jason、辰、立衡、w、zoro、刘岩、李森、岳溟、申信、changyu、王会迪、玄翼、张薇薇、马林、王孝诚

议题：

1. Github Project 需求功能进展同步
2. Milestone v1.1 release 计划梳理
3. 干扰检测中使用eBPF收集CPU调度延迟 proposal 介绍&讨论 （[#859](https://github.com/koordinator-sh/koordinator/pull/859)、[#860](https://github.com/koordinator-sh/koordinator/pull/860)）（delay next meeting）
4. CPU 抑制策略有 cpuset 和 cfs 两种，阿里是否在线上对比过其效果和差异？是否分别有适用的场景？
5. Koordinator Manager Recommender([https://koordinator.sh/docs/architecture/overview](https://koordinator.sh/docs/architecture/overview))资源画像这块有计划了吗？像 JVM 普遍给多少内存都能几乎用满，这种怎么做画像？
6. Koordlet的可观测性后面是啥计划？ 对集群参数调优时， 依赖于历史数据， 比如当CPU Burst策略设置为cpuBurstOnly时， 该怎么决定cpuBurstPercent的值设置为多少合适。
7. 对Koordinator的不同模块的不同特性，期望能提供一些压测案例、方法、数据出来。
action：

1. 确认 milestone v1.1 文档更新点

---
## 2022/11/29 #1

主持人：酒祝

参会人：酒祝、吕风、逐灵、宋涛、Jason、桑铎、林苍、CH、乔普、童超、立衡、醒宝、刘岩、申信、张凯、刘明、魏博锴、靳日阳、栗克宇、李森、祺琛、谢建超、HerbertDuan、张超、岳溟

议题：

1. Github Project 需求功能同步、对齐
2. 网络 QoS
3. 业务混部干扰因子度量
4. 单机侧支持的驱逐场景
5. 运行时的 CPU/Mem 等资源的动态调整逻辑
6. [https://github.com/koordinator-sh/koordinator/pull/835](https://github.com/koordinator-sh/koordinator/pull/835)
7. （欢迎补充）
8. 对后续社区周会形式/内容的建议和讨论
    1. 形式：议题讨论、Proposal 介绍、功能 Demo 演示...
    action：

1. QoS CPU 精细化的介绍文档

---

