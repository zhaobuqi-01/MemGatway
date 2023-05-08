package utils

import (
	"reflect"
	"sort"
	"strings"
)

// IsEqual 用于比较两个结构体是否相等，适用于不同类型的结构体。
func IsEqual(a, b interface{}) bool {
	// 获取两个结构体的反射值
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	// 检查两个结构体的类型是否相同
	if valA.Type() != valB.Type() {
		return false
	}

	// 检查两个结构体的字段数是否相同
	if valA.NumField() != valB.NumField() {
		return false
	}

	// 遍历结构体的字段并比较
	for i := 0; i < valA.NumField(); i++ {
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)

		// 忽略不可导出的字段
		if !fieldA.CanInterface() || !fieldB.CanInterface() {
			continue
		}

		// 对于字符串类型的字段，特殊处理IP地址或主机名列表
		if fieldA.Kind() == reflect.String && (strings.Contains(fieldA.String(), ",") || strings.Contains(fieldB.String(), ",")) {
			if !CompareStringSlices(ConvertAndSortIPs(fieldA.String()), ConvertAndSortIPs(fieldB.String())) {
				return false
			}
			continue
		}

		// 对于其他类型的字段，直接比较
		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			return false
		}
	}

	return true
}

// ConvertAndSortIPs 将逗号分隔的IP地址/主机名字符串转换为已排序的字符串切片
func ConvertAndSortIPs(ips string) []string {
	ipList := strings.Split(ips, ",")
	sort.Strings(ipList)
	return ipList
}

// CompareStringSlices 比较两个字符串切片是否相等
func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
