package cm

import (
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource/framework"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource/plugins/batchresource"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource/plugins/cpunormalization"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource/plugins/gpudeviceresource"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource/plugins/midresource"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource/plugins/resourceamplification"
	"github.com/koordinator-sh/koordinator/pkg/webhook/cm/plugins/sloconfig"
	corev1 "k8s.io/api/core/v1"
)

func CreateCheckersAll(oldConfig *corev1.ConfigMap, config *corev1.ConfigMap, needUnmarshal bool) interface{} {
	return []interface{}{
		sloconfig.NewColocationConfigChecker(oldConfig, config, needUnmarshal),  // ✅ 从  cfg 中获取对应的  []configuration.NodeCfgProfile
		sloconfig.NewResourceThresholdChecker(oldConfig, config, needUnmarshal), // ✅ 从  cfg 中获取对应的  []configuration.ResourceThresholdCfg
		sloconfig.NewResourceQOSChecker(oldConfig, config, needUnmarshal),       // ✅ 从  cfg 中获取对应的  []configuration.ResourceQOSCfg
		sloconfig.NewSystemConfigChecker(oldConfig, config, needUnmarshal),      // ✅ 从  cfg 中获取对应的  []configuration.NodeCfgProfile
		sloconfig.NewCPUBurstChecker(oldConfig, config, needUnmarshal),          // ✅ 从  cfg 中获取对应的  []configuration.CPUBurstCfg
	}
}

var (
	// SetupPlugins implement the setup for node resource plugin.
	setupPlugins = []framework.SetupPlugin{ // 它为插件设置了控制器客户端、方案和事件记录器等选项。
		&cpunormalization.Plugin{},
		&batchresource.Plugin{},
		&gpudeviceresource.Plugin{},
	}
	_ = new(batchresource.Plugin).Setup    // ✅
	_ = new(cpunormalization.Plugin).Setup // ✅
	// NodePreUpdatePlugin implements node resource pre-updating.
	nodePreUpdatePlugins = []framework.NodePreUpdatePlugin{
		&batchresource.Plugin{},
	}
	_ = new(batchresource.Plugin).PreUpdate // ✅ 更新node 前 更新  nrt zone-level  batch cpu
	// NodePreparePlugin implements node resource preparing for the calculated results.
	nodePreparePlugins = []framework.NodePreparePlugin{
		&cpunormalization.Plugin{}, // should be first
		&resourceamplification.Plugin{},
		&midresource.Plugin{},
		&batchresource.Plugin{},
		&gpudeviceresource.Plugin{},
	}
	_ = new(cpunormalization.Plugin).Prepare // ✅
	_ = new(midresource.Plugin).Prepare      // ✅
	// NodeSyncPlugin implements the check of resource updating.
	nodeStatusCheckPlugins = []framework.NodeStatusCheckPlugin{
		&midresource.Plugin{},
		&batchresource.Plugin{},
		&gpudeviceresource.Plugin{},
	}
	_ = new(midresource.Plugin).NeedSync   // ✅
	_ = new(batchresource.Plugin).NeedSync // ✅
	// nodeMetaCheckPlugins implements the check of node meta updating.
	nodeMetaCheckPlugins = []framework.NodeMetaCheckPlugin{
		&cpunormalization.Plugin{},
		&resourceamplification.Plugin{},
		&gpudeviceresource.Plugin{},
	}
	// ResourceCalculatePlugin implements resource counting and overcommitment algorithms.
	_ = new(cpunormalization.Plugin).NeedSyncMeta // ✅ 是否更新节点声明中的 cpu 归一化比例
	// ResourceCalculatePlugin implements resource counting and overcommitment algorithms.
	// 它根据node和NodeMetric计算节点资源，并生成计算的节点资源列表节点资源项。
	// 所有节点资源项都将合并到' noderresource '中，作为方法的中间结果计算阶段。
	// 在NodeMetric异常的情况下，插件可以在其中实现降级计算阶段。
	// compute插件还负责在主机托管时重置相应的节点资源配置为禁用。
	resourceCalculatePlugins = []framework.ResourceCalculatePlugin{
		&cpunormalization.Plugin{},
		&resourceamplification.Plugin{},
		&midresource.Plugin{},
		&batchresource.Plugin{},
		&gpudeviceresource.Plugin{},
	}
	_ = new(cpunormalization.Plugin).Calculate // ✅ node topo 和 strategy --> 添加声明
	_ = new(midresource.Plugin).Calculate      // ✅ 矫正 NodeMetric 中的 可回收资源，并记录到 nr
	_ = new(batchresource.Plugin).Calculate    // ✅ 统计 batch pod 使用的资源, 以及各 numa 使用情况

	_ = new(cpunormalization.Plugin).Reset
	_ = new(midresource.Plugin).Reset
	_ = new(batchresource.Plugin).Reset
)
