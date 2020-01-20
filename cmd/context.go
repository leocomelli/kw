package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/leocomelli/kw/pkg/common"
	"github.com/leocomelli/kw/pkg/config"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	getExample = templates.Examples(`
		# List all contexts.
		kw ctx

		# List all contexts in ps output format with more information (such as namespace).
		kw ctx -o wide

		# Modify the current context using the interactive mode
		kw ctx -i

		# Modify the current context
		kw ctx minikube

		# Switch to the previous context
		kw ctx -

		# Modify the current context and its namespace
		kw ctx minikube:kube-system

		# Modify the current context and switch to the previous namespace
		kw ctx minikube:-
		`)
)

// ContextOptions contains the input to the get command.
type ContextOptions struct {
	Output         string
	NoHeaders      bool
	Interactive    bool
	Config         *clientcmdapi.Config
	PahtOptions    *clientcmd.PathOptions
	KubeWideConfig *config.KubeWideConfig

	genericclioptions.IOStreams
}

func newContextOptions(s genericclioptions.IOStreams) *ContextOptions {
	configAccess := clientcmd.NewDefaultPathOptions()
	c, err := configAccess.GetStartingConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(PreFlightExitCode)
	}

	kw, err := config.NewKubeWideConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(PreFlightExitCode)
	}

	return &ContextOptions{
		Config:         c,
		PahtOptions:    configAccess,
		KubeWideConfig: kw,
		IOStreams:      s,
	}
}

// NewCmdContext creates a command object for the context actions
func NewCmdContext(streams genericclioptions.IOStreams) *cobra.Command {
	o := newContextOptions(streams)

	cmd := &cobra.Command{
		Use:     "ctx",
		Aliases: []string{"c", "context"},
		Short:   "Manage the context and namespace",
		Args:    cobra.MaximumNArgs(1),
		Example: getExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := len(args)

			if l == 0 && !o.Interactive {
				o.list()
			} else {
				var (
					context, namespace string
					err                error
				)

				if o.Interactive {
					keys := make([]string, 0, len(o.Config.Contexts))
					for k := range o.Config.Contexts {
						keys = append(keys, k)
					}
					context, err = common.InteractiveMode(keys)
					if err != nil {
						return err
					}
				} else {
					context, namespace = o.parseContextArg(args[0])
				}

				return o.set(context, namespace)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&o.Interactive, "interactive", "i", o.Interactive, "Enable interactive mode.")
	cmd.Flags().BoolVar(&o.NoHeaders, "no-headers", o.NoHeaders, "Does not print the headers.")
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "Output in the plain-text format with any additional information.")

	return cmd
}

func (o *ContextOptions) parseContextArg(name string) (string, string) {
	if name == PreviousIdentifier {
		return o.KubeWideConfig.PreviousContext(), ""
	}

	if params := strings.Split(name, ":"); len(params) >= 2 {
		if params[1] == PreviousIdentifier {
			return params[0], o.KubeWideConfig.PreviousNamespace()
		}
		return params[0], params[1]
	}
	return name, ""
}

func (o *ContextOptions) isWide() bool {
	return o.Output == "wide"
}

func (o *ContextOptions) set(ctx, ns string) error {
	// Preserve information about the current context to write it
	// to the internal file. In the future, we can use this information
	// to define to previous context and namespace.
	previousContext, ok := o.Config.Contexts[o.Config.CurrentContext]
	if ok {
		if o.Config.CurrentContext != ctx {
			o.KubeWideConfig.SetPreviousContext(o.Config.CurrentContext)
		}

		if o.Config.Contexts[o.Config.CurrentContext].Namespace != ns {
			o.KubeWideConfig.SetPreviousNamespace(previousContext.Namespace)
		}
	}

	newContext, ok := o.Config.Contexts[ctx]
	if !ok {
		return fmt.Errorf("context not found: %s", ctx)
	}

	o.Config.CurrentContext = ctx
	if ns != "" {
		newContext.Namespace = ns
	}

	err := clientcmd.ModifyConfig(o.PahtOptions, *o.Config, true)
	if err != nil {
		return fmt.Errorf("error when modifying the current context: %w", err)
	}

	// We should write the previous info in the internal file only when
	// the context has been successfully modified
	err = o.KubeWideConfig.Write()
	if err != nil {
		return err
	}

	return nil
}

func (o *ContextOptions) list() {
	var headers []string
	if !o.NoHeaders {
		headers = []string{"  CONTEXT", "NAMESPACE"}
		if !o.isWide() {
			headers = headers[:1]
		}
	}

	var data [][]string
	for k, v := range o.Config.Contexts {
		current := " "
		if k == o.Config.CurrentContext {
			current = "*"
		}
		if o.isWide() {
			data = append(data, []string{fmt.Sprintf("%s %s", current, k), v.Namespace})
			continue
		}
		data = append(data, []string{fmt.Sprintf("%s %s", current, k)})
	}

	common.TabPrint(os.Stdout, headers, data)
}
