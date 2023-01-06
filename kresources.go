package kresources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"github.com/scottd018/kresources/pkg/clusters"
	"github.com/scottd018/kresources/pkg/manifests"
)

// FromFiles returns a list of Kubernetes unstructured resource objects from a set of given
// file paths.
func FromFiles(files ...string) ([]*unstructured.Unstructured, error) {
	manifestResources := []*unstructured.Unstructured{}

	// get the resources from each of the files and append it to the array
	// of resources
	for i := range files {
		// create a new instance of a manifest object
		manifest, err := manifests.NewManifest(files[i])
		if err != nil {
			return []*unstructured.Unstructured{}, fmt.Errorf(
				"unable to create instance of manifest object - %w",
				err,
			)
		}

		manifestResources = append(manifestResources, manifest.Resources...)
	}

	return manifestResources, nil
}

// FromClusterFromFiles returns a list of Kubernetes unstructured resource objects from a set of given file paths.
// Resources are fetched from the cluster to determine their current state in the cluster.
func FromClusterFromFiles(client dynamic.Interface, files ...string) ([]*unstructured.Unstructured, error) {
	fileResources, err := FromFiles(files...)
	if err != nil {
		return []*unstructured.Unstructured{}, fmt.Errorf(
			"unable to retrieve resources from files - %w",
			err,
		)
	}

	clusterResources := []*unstructured.Unstructured{}

	// loop through the file resources and retrieve the unstructured resources from the cluster
	for i := range fileResources {
		clusterResource, err := clusters.NewClusterResource(
			&clusters.ClusterResourceInput{
				Group:     fileResources[i].GroupVersionKind().Group,
				Version:   fileResources[i].GroupVersionKind().Version,
				Kind:      fileResources[i].GroupVersionKind().Kind,
				Name:      fileResources[i].GetName(),
				Namespace: fileResources[i].GetNamespace(),
				Client:    client,
			},
		)

		if err != nil {
			return []*unstructured.Unstructured{}, fmt.Errorf(
				"unable to create cluster resource object from file resource - %w",
				err,
			)
		}

		// retrieve the resource from the cluster
		if err := clusterResource.Read(); err != nil {
			return []*unstructured.Unstructured{}, fmt.Errorf(
				"unable to read cluster resource object from cluster - %w",
				err,
			)
		}

		clusterResources = append(clusterResources, clusterResource.Resource)
	}

	return clusterResources, nil
}
