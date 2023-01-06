package manifests

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// FileResource is a representation of a resource which is extracted from a file.
type FileResource struct {
	Resource *unstructured.Unstructured
}

// NewFileResource returns a new instance of a FileResource object.
func NewFileResource(object *unstructured.Unstructured) *FileResource {
	return &FileResource{
		Resource: object,
	}
}
