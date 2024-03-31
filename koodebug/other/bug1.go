package other

import (
	corev1 "k8s.io/api/core/v1"
	quotav1 "k8s.io/apiserver/pkg/quota/v1"
	"k8s.io/client-go/tools/cache"
)

func LessThanOrEqual(a corev1.ResourceList, b corev1.ResourceList) (bool, []corev1.ResourceName) {
	result := true
	var resourceNames []corev1.ResourceName
	for key, value := range b { // a独有的没有计算
		if other, found := a[key]; found {
			if other.Cmp(value) > 0 {
				result = false
				resourceNames = append(resourceNames, key)
			}
		}
	}
	return result, resourceNames
}
func LessThanOrEqualCompletely(a corev1.ResourceList, b corev1.ResourceList) bool {
	result := true
	delta := quotav1.Subtract(a, b)
	for _, value := range delta {
		if value.Value() > 0 {
			result = false
			break
		}
	}

	return result
}
func Subtract(a corev1.ResourceList, b corev1.ResourceList) corev1.ResourceList {
	result := corev1.ResourceList{}
	for key, value := range a {
		quantity := value.DeepCopy()
		if other, found := b[key]; found {
			quantity.Sub(other)
		}
		result[key] = quantity
	}
	for key, value := range b {
		if _, found := result[key]; !found {
			quantity := value.DeepCopy()
			quantity.Neg()
			result[key] = quantity
		}
	}
	return result
}

func OnPodDelete(obj interface{}) {
	var pod *corev1.Pod
	switch t := obj.(type) {
	case *corev1.Pod:
		pod = t
	case cache.DeletedFinalStateUnknown:
		pod, _ = t.Obj.(*corev1.Pod)
	}
	if pod == nil {
		return
	}
}
