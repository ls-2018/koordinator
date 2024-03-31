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

package helpers

import (
	"fmt"
	"time"

	"k8s.io/klog/v2"

	slov1alpha1 "github.com/koordinator-sh/koordinator/apis/slo/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metriccache"
)

var (
	timeNow = time.Now
)

func CollectorNodeMetricLast(metricCache over_metriccache.MetricCache, queryMeta over_metriccache.MetricMeta, metricCollectInterval time.Duration) (float64, error) {
	queryParam := GenerateQueryParamsLast(metricCollectInterval * 2)
	result, err := CollectNodeMetrics(metricCache, *queryParam.Start, *queryParam.End, queryMeta)
	if err != nil {
		return 0, err
	}
	return result.Value(queryParam.Aggregate)
}

func CollectNodeMetrics(metricCache over_metriccache.MetricCache, start, end time.Time, queryMeta over_metriccache.MetricMeta) (over_metriccache.AggregateResult, error) {
	querier, err := metricCache.Querier(start, end)
	if err != nil {
		return nil, err
	}

	aggregateResult := over_metriccache.DefaultAggregateResultFactory.New(queryMeta)
	if err := querier.Query(queryMeta, nil, aggregateResult); err != nil {
		return nil, err
	}
	return aggregateResult, nil
}

func CollectAllHostAppMetricsLast(hostApps []slov1alpha1.HostApplicationSpec, metricCache over_metriccache.MetricCache,
	metricResource over_metriccache.MetricResource, metricCollectInterval time.Duration) map[string]float64 {
	queryParam := GenerateQueryParamsLast(metricCollectInterval * 2)
	return CollectAllHostAppMetrics(hostApps, metricCache, *queryParam, metricResource)
}

func CollectAllHostAppMetrics(hostApps []slov1alpha1.HostApplicationSpec, metricCache over_metriccache.MetricCache,
	queryParam over_metriccache.QueryParam, metricResource over_metriccache.MetricResource) map[string]float64 {
	appsMetrics := make(map[string]float64)
	querier, err := metricCache.Querier(*queryParam.Start, *queryParam.End)
	if err != nil {
		klog.Warningf("build host application querier failed, error: %v", err)
		return appsMetrics
	}
	for _, hostApp := range hostApps {
		queryMeta, err := metricResource.BuildQueryMeta(over_metriccache.MetricPropertiesFunc.HostApplication(hostApp.Name))
		if err != nil || queryMeta == nil {
			klog.Warningf("build host application %s query meta failed, kind: %s, error: %v", hostApp.Name, queryMeta, err)
			continue
		}
		aggregateResult := over_metriccache.DefaultAggregateResultFactory.New(queryMeta)
		if err := querier.Query(queryMeta, nil, aggregateResult); err != nil {
			klog.Warningf("query host application %s metric failed, kind: %s, error: %v", hostApp.Name, queryMeta.GetKind(), err)
			continue
		}
		if aggregateResult.Count() == 0 {
			klog.V(5).Infof("query host application %s metric is empty, kind: %s", hostApp.Name, queryMeta.GetKind())
			continue
		}
		value, err := aggregateResult.Value(queryParam.Aggregate)
		if err != nil {
			klog.Warningf("aggregate host application %s metric failed, kind: %s, error: %v", hostApp.Name, queryMeta.GetKind(), err)
			continue
		}
		appsMetrics[hostApp.Name] = value
	}
	return appsMetrics
}

func CollectPodMetricLast(metricCache over_metriccache.MetricCache, queryMeta over_metriccache.MetricMeta,
	metricCollectInterval time.Duration) (float64, error) {
	queryParam := GenerateQueryParamsLast(metricCollectInterval * 2)
	result, err := CollectPodMetric(metricCache, queryMeta, *queryParam.Start, *queryParam.End)
	if err != nil {
		return 0, err
	}
	return result.Value(queryParam.Aggregate)
}

func CollectPodMetric(metricCache over_metriccache.MetricCache, queryMeta over_metriccache.MetricMeta, start, end time.Time) (over_metriccache.AggregateResult, error) {
	querier, err := metricCache.Querier(start, end)
	if err != nil {
		return nil, err
	}
	aggregateResult := over_metriccache.DefaultAggregateResultFactory.New(queryMeta)
	if err := querier.Query(queryMeta, nil, aggregateResult); err != nil {
		return nil, err
	}
	return aggregateResult, nil
}

func CollectContainerResMetricLast(metricCache over_metriccache.MetricCache, queryMeta over_metriccache.MetricMeta,
	metricCollectInterval time.Duration) (float64, error) {
	queryParam := GenerateQueryParamsLast(metricCollectInterval * 2)
	querier, err := metricCache.Querier(*queryParam.Start, *queryParam.End)
	if err != nil {
		return 0, err
	}
	aggregateResult := over_metriccache.DefaultAggregateResultFactory.New(queryMeta)
	if err := querier.Query(queryMeta, nil, aggregateResult); err != nil {
		return 0, err
	}
	return aggregateResult.Value(queryParam.Aggregate)
}

func CollectContainerThrottledMetric(metricCache over_metriccache.MetricCache, containerID *string,
	metricCollectInterval time.Duration) (over_metriccache.AggregateResult, error) {
	if containerID == nil {
		return nil, fmt.Errorf("container is nil")
	}

	queryEndTime := time.Now()
	queryStartTime := queryEndTime.Add(-metricCollectInterval)
	querier, err := metricCache.Querier(queryStartTime, queryEndTime)
	if err != nil {
		return nil, err
	}

	queryParam := over_metriccache.MetricPropertiesFunc.Container(*containerID)
	queryMeta, err := over_metriccache.ContainerCPUThrottledMetric.BuildQueryMeta(queryParam)
	if err != nil {
		return nil, err
	}

	aggregateResult := over_metriccache.DefaultAggregateResultFactory.New(queryMeta)
	if err := querier.Query(queryMeta, nil, aggregateResult); err != nil {
		return nil, err
	}

	return aggregateResult, nil
}

func GenerateQueryParamsAvg(windowDuration time.Duration) *over_metriccache.QueryParam {
	end := time.Now()
	start := end.Add(-windowDuration)
	queryParam := &over_metriccache.QueryParam{
		Aggregate: over_metriccache.AggregationTypeAVG,
		Start:     &start,
		End:       &end,
	}
	return queryParam
}

func GenerateQueryParamsLast(windowDuration time.Duration) *over_metriccache.QueryParam {
	end := time.Now()
	start := end.Add(-windowDuration)
	queryParam := &over_metriccache.QueryParam{
		Aggregate: over_metriccache.AggregationTypeLast,
		Start:     &start,
		End:       &end,
	}
	return queryParam
}

func Query(querier over_metriccache.Querier, resource over_metriccache.MetricResource, properties map[over_metriccache.MetricProperty]string) (over_metriccache.AggregateResult, error) {
	queryMeta, err := resource.BuildQueryMeta(properties)
	if err != nil {
		return nil, err
	}

	aggregateResult := over_metriccache.DefaultAggregateResultFactory.New(queryMeta)
	if err := querier.Query(queryMeta, nil, aggregateResult); err != nil {
		return nil, err
	}

	return aggregateResult, nil
}
