package main

import (
	"fmt"
	koordletutil "github.com/koordinator-sh/koordinator/pkg/koordlet/util"
	util "github.com/koordinator-sh/koordinator/pkg/koordlet/util"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/util/system"
	corev1 "k8s.io/api/core/v1"
)

func main() {

	fmt.Println(koordletutil.GetBECPUSetPathsByMaxDepth(koordletutil.PodCgroupPathRelativeDepth))
	fmt.Println(util.GetRootCgroupCPUSetDir(corev1.PodQOSBestEffort))
	fmt.Println(system.GetRootCgroupSubfsDir(system.CgroupCPUSetDir))
}
