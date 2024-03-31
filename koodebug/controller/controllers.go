package controller

import (
	"flag"
	"github.com/koordinator-sh/koordinator/cmd/koord-manager/options"
	"github.com/koordinator-sh/koordinator/pkg/quota-controller/profile"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/nodemetric"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/noderesource"
	"github.com/koordinator-sh/koordinator/pkg/slo-controller/nodeslo"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var controllerInitFlags = map[string]func(*flag.FlagSet){
	noderesource.Name: noderesource.InitFlags,
}

var controllerAddFuncs = map[string]func(manager.Manager) error{
	nodemetric.Name:   nodemetric.Add,
	noderesource.Name: noderesource.Add,
	nodeslo.Name:      nodeslo.Add, // ✅ 将  configmap 、node 配置对应 nodesslo
	profile.Name:      profile.Add, // ✅    elasticquotaprofiles 、 elasticquotas
}

func main() {
	opts := options.NewOptions()
	opts.InitFlags(flag.CommandLine) // 主要逻辑 🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄🦄
}
