/*
Copyright 2022 The Koordinator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package koordlet

import (
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/framework"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/blkio"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/cgreconcile"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/cpuburst"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/cpuevict"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/cpusuppress"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/memoryevict"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/resctrl"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/qosmanager/plugins/sysreconcile"
)

var (
	StrategyPlugins = map[string]framework.QOSStrategyFactory{ // 指标收集插件
		blkio.BlkIOReconcileName:               blkio.New,        // ✅ 调整对应设备的速度限制
		cgreconcile.CgroupReconcileName:        cgreconcile.New,  // ✅ 根据全局配置，设置， memory.low < memory.min < memory.high 以调整
		cpuburst.CPUBurstName:                  cpuburst.New,     // ✅ 根据全局配置、pod身上的，根据 sharePool使用率 、定时协调容器状态  动态调整 cfs_burst_us\cfs_quota_us			calcStaticCPUBurstVal、genOperationByContainer、limiter.Allow(令牌桶)
		cpuevict.CPUEvictName:                  cpuevict.New,     // ✅ batch-cpu  be_resource
		cpusuppress.CPUSuppressName:            cpusuppress.New,  // ✅ 通过cfs_quota_us、cpuset 对 be pod 进行压制
		memoryevict.MemoryEvictName:            memoryevict.New,  // ✅ node_memory_usage
		resctrl.ResctrlReconcileName:           resctrl.New,      // ✅ 将 cgroup 中的task id 加入到  resctrl group 中
		sysreconcile.SystemConfigReconcileName: sysreconcile.New, // ✅ 定时更新 /proc/sys/vm/min_free_kbytes、watermark_scale_factor  kernel/mm/memcg_reaper/reap_background
	}
)
