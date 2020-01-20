package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/leocomelli/kw/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestParseContextArg(t *testing.T) {

	f, _ := ioutil.TempFile("", "kw_config")
	defer os.Remove(f.Name())

	os.Setenv("KW_CONFIG", f.Name())

	previousContext, _ := config.NewKubeWideConfig()
	previousContext.SetPreviousContext("previousCtx")

	previousNamespace, _ := config.NewKubeWideConfig()
	previousNamespace.SetPreviousNamespace("previousNs")

	tests := []struct {
		TestName  string
		Opts      *ContextOptions
		Arg       string
		Context   string
		Namespace string
	}{
		{"only context", &ContextOptions{}, "gke_cluster", "gke_cluster", ""},
		{"context and namespace", &ContextOptions{}, "gke_cluster:kube-system", "gke_cluster", "kube-system"},
		{"context and empty namespace", &ContextOptions{}, "gke_cluster:", "gke_cluster", ""},
		{"previous context", &ContextOptions{KubeWideConfig: previousContext}, "-", "previousCtx", ""},
		{"previous namespace", &ContextOptions{KubeWideConfig: previousNamespace}, "gke_cluster:-", "gke_cluster", "previousNs"},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			ctx, ns := tt.Opts.parseContextArg(tt.Arg)
			assert.Equal(t, tt.Context, ctx)
			assert.Equal(t, tt.Namespace, ns)
		})
	}
}
