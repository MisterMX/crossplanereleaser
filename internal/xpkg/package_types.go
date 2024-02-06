package xpkg

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Package struct {
	PackageObjects []*unstructured.Unstructured
	MetaObjects    []*unstructured.Unstructured
	ExampleObjects []*unstructured.Unstructured
}
