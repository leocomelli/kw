package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

const (
	defaultPath = "~/.kube/.kw.yml"

	previousContextKey   = "context"
	previousNamespaceKey = "namespace"
)

// KubeWideConfig represents the internal data
type KubeWideConfig struct {
	pathname string            `yaml:"-"`
	Previous map[string]string `yaml:"previous"`
}

// NewKubeWideConfig creates a new internal configuration
// and its dependencies
func NewKubeWideConfig() (*KubeWideConfig, error) {
	cc := &KubeWideConfig{
		Previous: map[string]string{
			previousContextKey:   "",
			previousNamespaceKey: "",
		},
	}
	err := cc.path()
	if err != nil {
		return nil, err
	}

	err = cc.read()
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *KubeWideConfig) path() error {
	c.pathname = os.Getenv("KW_CONFIG")
	if c.pathname == "" {
		var err error
		c.pathname, err = homedir.Expand(defaultPath)
		if err != nil {
			return fmt.Errorf("could not define the default path: %w", err)
		}
	}

	return nil
}

func (c *KubeWideConfig) read() error {
	f, err := ioutil.ReadFile(c.pathname)
	if err != nil {
		return fmt.Errorf("error reading the kw config file: %w", err)
	}
	return yaml.Unmarshal(f, c)
}

// Write the internal config file to the specified path
func (c *KubeWideConfig) Write() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error writing the kw config file: %w", err)
	}

	return ioutil.WriteFile(c.pathname, b, 0644)
}

// PreviousContext returns the previous context, otherwise empty
func (c *KubeWideConfig) PreviousContext() string {
	if pc, ok := c.Previous[previousContextKey]; ok {
		return pc
	}
	return ""
}

// PreviousNamespace returns the previous namespace, otherwise empty
func (c *KubeWideConfig) PreviousNamespace() string {
	if pc, ok := c.Previous[previousNamespaceKey]; ok {
		return pc
	}
	return ""
}

// SetPreviousContext sets the previous context
func (c *KubeWideConfig) SetPreviousContext(name string) {
	c.Previous[previousContextKey] = name
}

// SetPreviousNamespace sets the previous namespace
func (c *KubeWideConfig) SetPreviousNamespace(name string) {
	c.Previous[previousNamespaceKey] = name
}
