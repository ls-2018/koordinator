package system

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println(NewCommonSystemResource(KidledRelativePath, KidledScanPeriodInSecondsFileName, GetSysRootDir).Path(""))
	fmt.Println(NewCommonSystemResource(KidledRelativePath, KidledUseHierarchyFileFileName, GetSysRootDir).Path(""))
	fmt.Println(KidledScanPeriodInSeconds.Path(""))
	fmt.Println(sumUint64Slice2([]uint64{1, 2, 3, 4, 5, 6}))
}
func sumUint64Slice2(nums ...[]uint64) uint64 {
	var total uint64
	for _, v := range nums {
		for i := kidledColdBoundary; i < len(v); i++ {
			total = total + v[i]
		}
	}
	return total
}
