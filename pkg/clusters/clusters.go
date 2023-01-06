package clusters

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// Cluster represents an object that is used to communicate with a Kubernetes cluster
// via a dynamic client.
type Cluster struct {
	Client dynamic.Interface
}

// // Resource represents a resource that exists in a Kubernetes cluster.
// type Resource interface {
// 	Context() context.Context
// 	ResourceClient() dynamic.ResourceInterface
// 	Unstructured() *unstructured.Unstructured
// 	GroupVersionResource() (*schema.GroupVersionResource, error)
// }

// ClusterAction represents REST functions related to a cluster resource.
type ClusterAction func(*ClusterResource) (*unstructured.Unstructured, error)

// NewCluster creates a new instance of a cluster.
func NewCluster(client dynamic.Interface) *Cluster {
	return &Cluster{Client: client}
}

// Create creates a resource in a cluster.
func (cluster *Cluster) Create(resource *ClusterResource) (*unstructured.Unstructured, error) {
	return resource.Client.Create(
		resource.Context(),
		resource.Resource,
		metav1.CreateOptions{},
	)
}

// Read reads a resource from a cluster.
func (cluster *Cluster) Read(resource *ClusterResource) (*unstructured.Unstructured, error) {
	return resource.Client.Get(
		resource.Context(),
		resource.Resource.GetName(),
		metav1.GetOptions{},
	)
}

// Update updates a resource in a cluster.
func (cluster *Cluster) Update(resource *ClusterResource) (*unstructured.Unstructured, error) {
	// we need to retrieve the current resource from the cluster because we need to set the resource
	// version on the resource for a proper update
	currentResource, err := cluster.Read(resource)
	if err != nil {
		return resource.Resource, fmt.Errorf("unable to get resource - %w", err)
	}

	resource.Resource.SetResourceVersion(currentResource.GetResourceVersion())

	return resource.Client.Update(
		resource.Context(),
		resource.Resource,
		metav1.UpdateOptions{},
	)
}

// Delete deletes a resource from a cluster.
func (cluster *Cluster) Delete(resource *ClusterResource) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{}, resource.Client.Delete(
		resource.Context(),
		resource.Resource.GetName(),
		metav1.DeleteOptions{},
	)
}
