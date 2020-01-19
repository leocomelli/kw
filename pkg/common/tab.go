package common

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const tab = '\t'

// TabPrint writes data to a specific writer using the tabular format
func TabPrint(w io.Writer, headers []string, data [][]string) {

	tw := &tabwriter.Writer{}
	tw.Init(w, 0, 8, 0, tab, 0)

	if len(headers) > 0 {
		fmt.Fprintln(tw, strings.Join(headers, string(tab)))
	}

	for _, d := range data {
		fmt.Fprintln(tw, strings.Join(d, string(tab)))
	}

	tw.Flush()
}
