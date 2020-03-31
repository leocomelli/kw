package cmd

import (
	"fmt"

	"github.com/leocomelli/kw/pkg/common"
	"github.com/leocomelli/kw/pkg/kubernetes"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	logExamples = templates.Examples(`
		# Streams logs from all pods in the namespace.
		kw logs -n kube-system

		# Streams logs from the pod specified in the namespace.
		kw logs -n kube-system -p kube-dns-5c446b66bd-p7s2f

		# Streams logs from an specific container in a given namespace and pod.
		kw ns cert-manager
		`)
)

// LogOptions contains the input to the get command.
type LogOptions struct {
	Namespace string
	Pod       string
	Container string
	NoColor   bool

	genericclioptions.IOStreams
}

// NewCmdLogs creates a command object for the logs actions
func NewCmdLogs(streams genericclioptions.IOStreams) *cobra.Command {
	o := &LogOptions{IOStreams: streams}

	cmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"l", "log"},
		Short:   "Streams logs from all containers of all matched pods",
		Args:    cobra.MaximumNArgs(1),
		Example: logExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Namespace == "" {
				return fmt.Errorf("namespace is required")
			}

			k, err := kubernetes.NewKubernetes()
			if err != nil {
				return err
			}

			stream := make(chan *kubernetes.LogStream)
			go k.Logs(o.Namespace, o.Pod, o.Container, stream)

			pc := common.NewPrintColor()

			var lastKey string
			colors := make(map[string]common.PrintFn)

			// initial padding, it will be changed based on the key length
			padding := 35
			for l := range stream {
				currentKey := fmt.Sprintf("%s/%s", l.Pod, l.Container)

				if o.NoColor {
					colors[currentKey] = pc.GetNoColor()
				} else {
					_, ok := colors[currentKey]
					if !ok {
						colors[currentKey] = pc.Get()
					}
				}

				if lastKey != currentKey {
					lastKey = currentKey
					if len(lastKey) > padding {
						padding = len(lastKey)
					}
				}

				tmpl := fmt.Sprintf("%%-%dv | %%s\n", padding)
				fmt.Printf(tmpl, colors[currentKey](fmt.Sprintf("%s", lastKey)), l.Message)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.NoColor, "no-color", o.NoColor, "disable ansi color output")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "match pods in the given namespace")
	cmd.Flags().StringVarP(&o.Pod, "pod", "p", o.Pod, "match pods by name")
	cmd.Flags().StringVarP(&o.Container, "container", "c", o.Container, "restrict which containers logs are shown for")

	return cmd
}
