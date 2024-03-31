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

package over_metricsadvisor

import (
	"time"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metriccache"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metricsadvisor/over_framework"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_statesinformer"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/resourceexecutor"
)

type MetricAdvisor interface {
	Run(stopCh <-chan struct{}) error
	HasSynced() bool
}

type metricAdvisor struct {
	options *over_framework.Options
	context *over_framework.Context
}

func NewMetricAdvisor(cfg *over_framework.Config, statesInformer over_statesinformer.StatesInformer, metricCache over_metriccache.MetricCache) MetricAdvisor {
	opt := &over_framework.Options{
		Config:         cfg,
		StatesInformer: statesInformer,
		MetricCache:    metricCache,
		CgroupReader:   resourceexecutor.NewCgroupReader(),
		PodFilters:     podFilters,
	}
	ctx := &over_framework.Context{
		DeviceCollectors: make(map[string]over_framework.DeviceCollector, len(devicePlugins)),
		Collectors:       make(map[string]over_framework.Collector, len(collectorPlugins)),
		State:            over_framework.NewSharedState(),
	}
	for name, device := range devicePlugins {
		ctx.DeviceCollectors[name] = device(opt)
	}
	for name, collector := range collectorPlugins {
		ctx.Collectors[name] = collector(opt)
	}

	c := &metricAdvisor{
		options: opt,
		context: ctx,
	}
	return c
}

func (m *metricAdvisor) HasSynced() bool {
	return over_framework.CollectorsHasStarted(m.context.Collectors)
}

func (m *metricAdvisor) Run(stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	if m.options.Config.CollectResUsedInterval < time.Second {
		klog.Infof("CollectResUsedInterval is %v, metric collector is disabled", m.options.Config.CollectResUsedInterval)
		return nil
	}

	defer m.shutdown()
	m.setup()

	defer klog.Info("shutting down metric advisor")
	klog.Info("Starting collector for NodeMetric")

	for name, dc := range m.context.DeviceCollectors {
		klog.V(4).Infof("ready to start device collector %v", name)
		if !dc.Enabled() {
			klog.V(4).Infof("device collector %v is not enabled, skip running", name)
			continue
		}
		go dc.Run(stopCh)
		klog.V(4).Infof("device collector %v start", name)
	}

	for name, collector := range m.context.Collectors {
		klog.V(4).Infof("ready to start collector %v", name)
		if !collector.Enabled() {
			klog.V(4).Infof("collector %v is not enabled, skip running", name)
			continue
		}
		go collector.Run(stopCh)
		klog.V(4).Infof("collector %v start", name)
	}

	klog.Info("Starting successfully")
	<-stopCh
	return nil
}

func (m *metricAdvisor) setup() {
	for name, dc := range m.context.DeviceCollectors {
		if !dc.Enabled() {
			klog.V(4).Infof("device collector %v is not enabled, skip setup", name)
			continue
		}
		dc.Setup(m.context)
	}
	for name, collector := range m.context.Collectors {
		if !collector.Enabled() {
			klog.V(4).Infof("collector %v is not enabled, skip setup", name)
			continue
		}
		collector.Setup(m.context)
	}
}

func (m *metricAdvisor) shutdown() {
	for name, dc := range m.context.DeviceCollectors {
		if !dc.Enabled() {
			klog.V(4).Infof("device collector %v is not enabled, skip shutdown", name)
			continue
		}
		dc.Shutdown()
	}
}
