package manifests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/scottd018/kresources/pkg/resources"
)

var (
	ErrFileMissing = errors.New("unable to locate file")
	ErrFileRead    = errors.New("unable to read file")

	ErrConvertUnstructured = errors.New("unable to convert to unstructured object")
)

// Manifest represents a Kubernetes file manifest.
type Manifest struct {
	File      string
	Content   []byte
	YAML      []*yaml.Node
	Resources []*unstructured.Unstructured
}

// NewManifest creates a new instance of a Manifest object.
func NewManifest(path string) (*Manifest, error) {
	// ensure the file path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file missing %s - %w", path, ErrFileMissing)
	}

	// read the contents of the file path
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("file error %s - %w", path, ErrFileRead)
	}

	// create the manifest object in memory
	manifest := &Manifest{
		File:    path,
		Content: content,
	}

	// read the yaml from the contents
	yamlContent, err := manifest.ToYAML()
	if err != nil {
		return manifest, fmt.Errorf("error retrieving yaml objects from path %s - %w", manifest.File, err)
	}

	manifest.YAML = yamlContent

	// loop through the yaml objects and create a new resource object from each one
	for i := range manifest.YAML {
		object, err := resources.ToUnstructured(manifest.YAML[i])
		if err != nil {
			return manifest, fmt.Errorf("%s - %w", ErrConvertUnstructured, err)
		}

		manifest.Resources = append(manifest.Resources, object)
	}

	return manifest, nil
}

// ToYAML returns a manifest as an array of YAML objects.
func (manifest *Manifest) ToYAML() ([]*yaml.Node, error) {
	yamlObjects := []*yaml.Node{}

	// create a decoder object to decode yaml objects from bytes
	decoder := yaml.NewDecoder(bytes.NewReader(manifest.Content))

	// loop through each document and add it to the array of yamlObjects
	for {
		// read the next document from the file
		var node yaml.Node

		err := decoder.Decode(node)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error decoding manifest %s - %w", manifest.File, err)
		}

		// append the document to the slice
		yamlObjects = append(yamlObjects, &node)
	}

	return yamlObjects, nil
}
