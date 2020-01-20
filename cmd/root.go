package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	// PreFlightExitCode defines exit the code for preflight checks
	PreFlightExitCode = 2
	// PreviousIdentifier defines that the previous value should be used
	PreviousIdentifier = "-"
)

// NewCmdKubeWide creates the `kw` command and its nested children.
func NewCmdKubeWide(in io.Reader, out, err io.Writer) *cobra.Command {
	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	cmds := &cobra.Command{
		Use:   "kw",
		Short: "kw is an extension of kubectl to help us manage our kubernetes clusters",
		Long:  ``,
	}

	cmds.AddCommand(NewCmdContext(ioStreams))
	cmds.AddCommand(NewCmdKubectl(ioStreams))
	cmds.AddCommand(NewCmdNamespace(ioStreams))

	return cmds
}
