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

package over_nodestorageinfo

import (
	"time"

	"go.uber.org/atomic"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	"github.com/koordinator-sh/koordinator/pkg/features"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metriccache"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metrics"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/over_framework"
	koordletutil "github.com/koordinator-sh/koordinator/pkg/koordlet/util"
)

const (
	CollectorName = "NodeStorageInfoCollector"
)

type nodeInfoCollector struct {
	collectInterval time.Duration
	storage         over_metriccache.KVStorage
	started         *atomic.Bool
}

func New(opt *over_framework.Options) over_framework.Collector {
	return &nodeInfoCollector{
		collectInterval: opt.Config.CollectNodeStorageInfoInterval,
		storage:         opt.MetricCache,
		started:         atomic.NewBool(false),
	}
}

func (n *nodeInfoCollector) Enabled() bool {
	return features.DefaultKoordletFeatureGate.Enabled(features.BlkIOReconcile)
}

func (n *nodeInfoCollector) Setup(s *over_framework.Context) {}

func (n *nodeInfoCollector) Run(stopCh <-chan struct{}) {
	go wait.Until(n.collectNodeLocalStorageInfo, n.collectInterval, stopCh)
}

func (n *nodeInfoCollector) Started() bool {
	return n.started.Load()
}

func (n *nodeInfoCollector) collectNodeLocalStorageInfo() {
	klog.V(6).Info("start collect node local storage info")

	localStorageInfo, err := koordletutil.GetLocalStorageInfo()
	if err != nil {
		klog.Warningf("failed to collect node local storage info, err: %s", err)
		over_metrics.RecordCollectNodeLocalStorageInfoStatus(err)
		return
	}

	nodeLocalStorageInfo := &over_metriccache.NodeLocalStorageInfo{}
	nodeLocalStorageInfo.DiskNumberMap = localStorageInfo.DiskNumberMap
	nodeLocalStorageInfo.NumberDiskMap = localStorageInfo.NumberDiskMap
	nodeLocalStorageInfo.PartitionDiskMap = localStorageInfo.PartitionDiskMap
	nodeLocalStorageInfo.VGDiskMap = localStorageInfo.VGDiskMap
	nodeLocalStorageInfo.LVMapperVGMap = localStorageInfo.LVMapperVGMap
	nodeLocalStorageInfo.MPDiskMap = localStorageInfo.MPDiskMap

	klog.V(6).Infof("collect node local storage info finished, nodeCPUInfo %v", localStorageInfo)
	n.storage.Set(over_metriccache.NodeLocalStorageInfoKey, nodeLocalStorageInfo)
	n.started.Store(true)
	over_metrics.RecordCollectNodeLocalStorageInfoStatus(nil)
}
