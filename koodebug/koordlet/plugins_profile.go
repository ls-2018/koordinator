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
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/beresource"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/coldmemoryresource"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/over_hostapplication"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/over_nodeinfo"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/over_noderesource"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/over_nodestorageinfo"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/pagecache"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/performance"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/podresource"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/podthrottled"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/collectors/sysresource"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/over_devices/gpu"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/over_framework"
)

// NOTE: map variables in this file can be overwritten for extension

var (
	devicePlugins = map[string]over_framework.DeviceFactory{
		gpu.DeviceCollectorName: gpu.New, // ✅ 采集 gpu pod、level、node 指标
	}

	collectorPlugins = map[string]over_framework.CollectorFactory{ // 指标收集
		over_noderesource.CollectorName:    over_noderesource.New,    // ✅
		beresource.CollectorName:           beresource.New,           // ✅
		over_nodeinfo.CollectorName:        over_nodeinfo.New,        // ✅
		over_nodestorageinfo.CollectorName: over_nodestorageinfo.New, // ✅
		podresource.CollectorName:          podresource.New,          // ✅
		podthrottled.CollectorName:         podthrottled.New,         // ✅
		performance.CollectorName:          performance.New,          // ✅ 最重要的
		sysresource.CollectorName:          sysresource.New,          // ✅ 依赖  noderesource 、podresource   节点使用-pod使用-hostApp使用=系统使用
		coldmemoryresource.CollectorName:   coldmemoryresource.New,   // ✅
		pagecache.CollectorName:            pagecache.New,            // ✅ 依赖 podFilters
		over_hostapplication.CollectorName: over_hostapplication.New, // ✅
	}

	podFilters = map[string]over_framework.PodFilter{
		podresource.CollectorName:  over_framework.DefaultPodFilter, // 返回 true 的 都不处理
		podthrottled.CollectorName: over_framework.DefaultPodFilter,
	}
)
