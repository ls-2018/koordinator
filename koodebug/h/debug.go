package main

import (
	"fmt"
	slov1alpha1 "github.com/koordinator-sh/koordinator/apis/slo/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/util/histogram"
	"sort"
)

var (
	// MinSampleWeight is the minimal weight of any sample (prior to including decaying factor)
	MinSampleWeight = 0.1
	// epsilon is the minimal weight kept in histograms, it should be small enough that old samples
	// (just inside MemoryAggregationWindowLength) added with MinSampleWeight are still kept
	epsilon = 0.001 * MinSampleWeight
	// DefaultHistogramBucketSizeGrowth is the default value for histogramBucketSizeGrowth.
	DefaultHistogramBucketSizeGrowth = 0.05
)

func main2() {
	options, err := histogram.NewExponentialHistogramOptions(1024, 0.025, 1.+DefaultHistogramBucketSizeGrowth, epsilon)

	fmt.Println(options, err)

	options, err = histogram.NewExponentialHistogramOptions(1<<31, 5<<20, 1.+DefaultHistogramBucketSizeGrowth, epsilon)
	fmt.Println(options, err)

	pairs := []int{1, 25, 123, 12, 31, 2}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i] < pairs[j]
	})
	fmt.Println(pairs)
}
func main() {

	var cfg *slov1alpha1.ResourceQOS
	fmt.Println(cfg.DeepCopy() == nil)

}
