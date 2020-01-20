package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/leocomelli/kw/pkg/common"
	"github.com/leocomelli/kw/pkg/config"
	"github.com/leocomelli/kw/pkg/kubernetes"
	"github.com/spf13/cobra"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	nsExamples = templates.Examples(`
		# List all namespaces.
		kw ns

		# Modify the current namespace using the interactive mode
		kw ns -i

		# Modify the current namespace
		kw ns cert-manager

		# Switch to the previous namespace
		kw ns -
		`)
)

// NamespaceOptions contains the input to the get command.
type NamespaceOptions struct {
	NoHeaders      bool
	Interactive    bool
	Config         *clientcmdapi.Config
	PahtOptions    *clientcmd.PathOptions
	KubeWideConfig *config.KubeWideConfig
	Kubernetes     *kubernetes.Kubernetes

	genericclioptions.IOStreams
}

func newNamespaceOptions(s genericclioptions.IOStreams) *NamespaceOptions {
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

	k, err := kubernetes.NewKubernetes()
	if err != nil {
		fmt.Println(err)
		os.Exit(PreFlightExitCode)
	}

	return &NamespaceOptions{
		Config:         c,
		PahtOptions:    configAccess,
		KubeWideConfig: kw,
		Kubernetes:     k,
		IOStreams:      s,
	}
}

// NewCmdNamespace creates a command object for the namespace actions
func NewCmdNamespace(streams genericclioptions.IOStreams) *cobra.Command {
	o := newNamespaceOptions(streams)

	cmd := &cobra.Command{
		Use:     "ns",
		Aliases: []string{"n", "namespace"},
		Short:   "Manage the namespaces",
		Args:    cobra.MaximumNArgs(1),
		Example: nsExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := len(args)

			if l == 0 && !o.Interactive {
				err := o.list()
				if err != nil {
					return err
				}
			} else {
				var namespace string

				if o.Interactive {
					ns, err := o.Kubernetes.Namespaces()
					if err != nil {
						return err
					}

					keys := make([]string, 0, len(ns))
					for _, n := range ns {
						keys = append(keys, n.GetName())
					}
					namespace, err = common.InteractiveMode(keys)
					if err != nil {
						return err
					}
				} else {
					namespace = args[0]
				}

				return o.set(namespace)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&o.Interactive, "interactive", "i", o.Interactive, "Enable interactive mode.")
	cmd.Flags().BoolVar(&o.NoHeaders, "no-headers", o.NoHeaders, "Does not print the headers.")

	return cmd
}

func (o *NamespaceOptions) set(ns string) error {
	if ns == PreviousIdentifier {
		ns = o.KubeWideConfig.PreviousNamespace()
	}

	// Preserve information about the current context to write it
	// to the internal file. In the future, we can use this information
	// to define to previous context and namespace.
	context, ok := o.Config.Contexts[o.Config.CurrentContext]
	if ok {
		if context.Namespace != ns {
			o.KubeWideConfig.SetPreviousNamespace(context.Namespace)
		}
	}

	context.Namespace = ns

	err := clientcmd.ModifyConfig(o.PahtOptions, *o.Config, true)
	if err != nil {
		return fmt.Errorf("error when modifying the current namespace: %w", err)
	}

	// We should write the previous info in the internal file only when
	// the context has been successfully modified
	err = o.KubeWideConfig.Write()
	if err != nil {
		return err
	}

	return nil
}

func (o *NamespaceOptions) list() error {
	var headers []string
	if !o.NoHeaders {
		headers = []string{"  NAME", "STATUS", "AGE"}
	}

	ns, err := o.Kubernetes.Namespaces()
	if err != nil {
		return err
	}

	var curNamespace string
	if ctx, ok := o.Config.Contexts[o.Config.CurrentContext]; ok {
		curNamespace = ctx.Namespace
	} // else current context does not exist

	var data [][]string
	for _, n := range ns {
		current := " "
		if n.GetName() == curNamespace {
			current = "*"
		}
		data = append(data, []string{fmt.Sprintf("%s %s", current, n.GetName()), n.Status.String(), translateTimestampSince(n.GetCreationTimestamp())})
	}

	common.TabPrint(os.Stdout, headers, data)

	return nil
}

// translateTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestampSince(timestamp meta.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}
