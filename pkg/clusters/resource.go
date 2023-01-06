package clusters

import (
	"context"
	"errors"
	"fmt"

	"github.com/scottd018/kresources/pkg/resources"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	ErrMissingGroupInput   = errors.New("missing Group field on ClusterResourceInput object")
	ErrMissingVersionInput = errors.New("missing Version field on ClusterResourceInput object")
	ErrMissingKindInput    = errors.New("missing Kind field on ClusterResourceInput object")
	ErrMissingNameInput    = errors.New("missing Name field on ClusterResourceInput object")
	ErrMissingClientInput  = errors.New("missing dynamic.Interface object as Client field on ClusterResourceInput object")
	ErrMissingResource     = errors.New("missing Resource object on ClusterResource object")
)

// ClusterResource is a representation of a resource which is extracted from a cluster.
type ClusterResource struct {
	Resource *unstructured.Unstructured
	Client   dynamic.ResourceInterface
	Input    *ClusterResourceInput
	Cluster  *Cluster
}

// ClusterResourceInput is a representation of the data required to create a cluster resource.
type ClusterResourceInput struct {
	Group     string
	Version   string
	Kind      string
	Name      string
	Namespace string

	Client  dynamic.Interface
	Context context.Context
}

// Validate validates a ClusterResourceInput.
func (input *ClusterResourceInput) Validate() error {
	// ensure group is set
	if input.Group == "" {
		return ErrMissingGroupInput
	}

	// ensure version is set
	if input.Version == "" {
		return ErrMissingVersionInput
	}

	// ensure kind is set
	if input.Kind == "" {
		return ErrMissingKindInput
	}

	// ensure name is set
	if input.Name == "" {
		return ErrMissingKindInput
	}

	// ensure client is set
	if input.Client == nil {
		return ErrMissingClientInput
	}

	return nil
}

// NewClusterResource creates a new instance of a cluster resource.
func NewClusterResource(input *ClusterResourceInput) (*ClusterResource, error) {
	// validate the input
	if err := input.Validate(); err != nil {
		return &ClusterResource{}, err
	}

	// set the context if one was not provided
	if input.Context == nil {
		input.Context = context.Background()
	}

	// create a base object that we will store for the unstructured object
	object := &unstructured.Unstructured{}
	object.SetGroupVersionKind(
		schema.GroupVersionKind{
			Group:   input.Group,
			Version: input.Version,
			Kind:    input.Kind,
		},
	)
	object.SetName(input.Name)
	object.SetNamespace(input.Namespace)

	// create the resource client
	gvr, err := resources.ToGroupVersionResource(object)
	if err != nil {
		return nil, fmt.Errorf("unable to get schema.GroupVersionResource - %w", err)
	}

	// create a new instance of a cluster resource object and return it
	clusterResource := &ClusterResource{
		Input:    input,
		Cluster:  NewCluster(input.Client),
		Resource: object,
		Client:   input.Client.Resource(*gvr).Namespace(object.GetNamespace()),
	}

	return clusterResource, nil
}

// Context returns the context from the cluster resource.
func (resource *ClusterResource) Context() context.Context {
	return resource.Input.Context
}

// GroupVersionResource returns the GVR for an object.
func (resource *ClusterResource) GroupVersionResource() (*schema.GroupVersionResource, error) {
	return resources.ToGroupVersionResource(resource.Resource)
}

// Create creates a cluster resource.
func (resource *ClusterResource) Create() error {
	// ensure the resource is set
	if resource.Resource == nil {
		return ErrMissingResource
	}

	return resource.perform(resource.Cluster.Create)
}

// Read reads a cluster resource from a cluster.
func (resource *ClusterResource) Read() error {
	return resource.perform(resource.Cluster.Read)
}

// Update updates a cluster resource.
func (resource *ClusterResource) Update() error {
	// ensure the resource is set
	if resource.Resource == nil {
		return ErrMissingResource
	}

	return resource.perform(resource.Cluster.Update)
}

// Delete deletes a cluster resource.
func (resource *ClusterResource) Delete() error {
	// ensure the resource is set
	if resource.Resource == nil {
		return ErrMissingResource
	}

	return resource.perform(resource.Cluster.Delete)
}

// perform is an entrypoint to perform a set of actions against a cluster.
func (resource *ClusterResource) perform(performAction ClusterAction) error {
	// run the cluster action which was passed in
	object, err := performAction(resource)
	if err != nil {
		return err
	}

	// set the retrieved object on the resource
	resource.Resource = object

	return nil
}
