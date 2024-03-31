package descheduler

import (
	"github.com/koordinator-sh/koordinator/pkg/descheduler/controllers/migration"
	myframework "github.com/koordinator-sh/koordinator/pkg/descheduler/framework"
	"github.com/koordinator-sh/koordinator/pkg/descheduler/framework/plugins/kubernetes/defaultevictor"
	"github.com/koordinator-sh/koordinator/pkg/descheduler/framework/plugins/loadaware"
	"os"
	"sigs.k8s.io/descheduler/pkg/framework"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/nodeutilization"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/podlifetime"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removeduplicates"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removefailedpods"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removepodshavingtoomanyrestarts"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removepodsviolatinginterpodantiaffinity"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removepodsviolatingnodeaffinity"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removepodsviolatingnodetaints"
	"sigs.k8s.io/descheduler/pkg/framework/plugins/removepodsviolatingtopologyspreadconstraint"
)

func Init() {
	os.Args = append(os.Args,
		"--config=/Users/acejilam/Desktop/work/koordinator/debug/descheduler/descheduler.yaml",
	)
}

func main() {
	//pluginFactory(args, f)

	{
		_ = defaultevictor.New // ✅
	}

	{
		_ = migration.New
		_ = migration.Reconciler{}.Start // 是一个 controller
		// watch PodMigrationJob、Reservation
		// add   ->waitingCollection
		// delete->arbitratedPodMigrationJobs
		// while :doOnceArbitrate
		// 更新 PodMigrationJob 状态、清理过期的 job和Reservation
	}
	{
		_ = podlifetime.New // 存活超过一定时间，就驱逐
		_ = podlifetime.PodLifeTime{}
	}

	{
		_ = loadaware.NewLowNodeLoad // 依赖 slo   koord
		_ = loadaware.LowNodeLoad{}
	}

	{
		_ = removefailedpods.New // 驱逐 failed状态的 pod
		_ = removefailedpods.RemoveFailedPods{}
	}

	{
		_ = removeduplicates.New // 移除重复 pod,尽可能均匀分布
		_ = removeduplicates.RemoveDuplicates{}

	}
	{
		_ = removepodshavingtoomanyrestarts.New // 移除 重启次数过多的
		_ = removepodshavingtoomanyrestarts.RemovePodsHavingTooManyRestarts{}
	}

	{
		_ = removepodsviolatingnodetaints.New // 驱逐节点上违反NoSchedule tts的pod
	}
	{
		_ = removepodsviolatinginterpodantiaffinity.New // 检查节点上是否存在当前pod不能容忍的其他pod

	}
	_ = nodeutilization.NewHighNodeUtilization // 从未被使用的节点中驱逐pod
	_ = nodeutilization.NewLowNodeUtilization  // 驱逐高使用率节点的pod

	_ = removepodsviolatingnodeaffinity.New             // 清除节点上违反节点亲和性的pod
	_ = removepodsviolatingtopologyspreadconstraint.New // 驱逐违反其拓扑扩展约束的pod

	var _ framework.DeschedulePlugin = &podlifetime.PodLifeTime{}

	var _ myframework.BalancePlugin = &loadaware.LowNodeLoad{}
	var _ framework.BalancePlugin = &nodeutilization.HighNodeUtilization{}
	var _ framework.BalancePlugin = &nodeutilization.LowNodeUtilization{}
	var _ framework.BalancePlugin = &removeduplicates.RemoveDuplicates{}
	var _ framework.BalancePlugin = &removepodsviolatingtopologyspreadconstraint.RemovePodsViolatingTopologySpreadConstraint{}

	var _ framework.DeschedulePlugin = &removefailedpods.RemoveFailedPods{}
	var _ framework.DeschedulePlugin = &removepodshavingtoomanyrestarts.RemovePodsHavingTooManyRestarts{}
	var _ framework.DeschedulePlugin = &removepodsviolatinginterpodantiaffinity.RemovePodsViolatingInterPodAntiAffinity{}
	var _ framework.DeschedulePlugin = &removepodsviolatingnodeaffinity.RemovePodsViolatingNodeAffinity{}
	var _ framework.DeschedulePlugin = &removepodsviolatingnodetaints.RemovePodsViolatingNodeTaints{}

}
