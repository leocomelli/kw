package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/cmd"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// NewCmdKubectl creates a command object that wraps the kubectl oficial command
func NewCmdKubectl(streams genericclioptions.IOStreams) *cobra.Command {
	cmdKubectlWrap := cmd.NewKubectlCommand(streams.In, streams.Out, streams.ErrOut)
	cmdKubectlWrap.Use = "ctl"
	cmdKubectlWrap.Aliases = []string{"control", "kubectl"}
	cmdKubectlWrap.Short = "Wraps the official kubectl command"

	return cmdKubectlWrap
}
