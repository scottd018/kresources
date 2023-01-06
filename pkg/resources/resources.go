package resources

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	resourceutil "sigs.k8s.io/kubebuilder/v3/pkg/model/resource"
)

var (
	ErrConvertUnstructured = errors.New("unable to convert to unstructured object")
	ErrMissingUnstructured = errors.New("unable to find unstructured object")
	ErrMissingKind         = errors.New("unable to find kind")
)

// ToGroupVersionResource returns a GroupVersionResource object given an unstructured object.
func ToGroupVersionResource(object *unstructured.Unstructured) (*schema.GroupVersionResource, error) {
	// return immediately if we happen to have passed in a nil object
	if object == nil {
		return &schema.GroupVersionResource{}, ErrMissingUnstructured
	}

	// retrieve the plural version of the resource name
	resourceName, err := PluralResourceName(object.GroupVersionKind())
	if err != nil {
		return &schema.GroupVersionResource{}, fmt.Errorf("unable to retrieve resource name - %w", err)
	}

	// create the gvr object and return.  we need to pluralize the object in the manner in which
	// we interact with a cluster via kubectl (e.g. Kind=Pod, Resource=pods, kubectl get pods)
	gvr := &schema.GroupVersionResource{
		Group:    object.GroupVersionKind().Group,
		Version:  object.GroupVersionKind().Version,
		Resource: resourceName,
	}

	return gvr, nil
}

// PluralResourceName returns the plural resource name of a GVK object.
func PluralResourceName(gvk schema.GroupVersionKind) (string, error) {
	if gvk.Kind == "" {
		return "", fmt.Errorf("unable to pluralize resource - %w", ErrMissingKind)
	}

	return strings.ToLower(resourceutil.RegularPlural(gvk.Kind)), nil
}

// ToUnstructured converts a variety of known objects into Kubernetes unstructured objects.  It returns
// an error with an unsupported message in the instance that the input type is unknown.
func ToUnstructured(input interface{}) (*unstructured.Unstructured, error) {
	switch object := input.(type) {
	case yaml.Node:
		return FromYAML(&object)
	case *yaml.Node:
		return FromYAML(object)
	}

	return &unstructured.Unstructured{}, fmt.Errorf("unsupported type [%T] - %w", input, ErrConvertUnstructured)
}

// FromYAML converts a YAML node object into a Kubernetes unstructured object.
func FromYAML(yamlNode *yaml.Node) (*unstructured.Unstructured, error) {
	// marshal the YAML node object into JSON
	jsonBytes, err := yaml.Marshal(yamlNode)
	if err != nil {
		return nil, err
	}

	// create a new unstructured object and set the JSON data
	object := &unstructured.Unstructured{}
	if err := object.UnmarshalJSON(jsonBytes); err != nil {
		return nil, err
	}

	return object, nil
}
