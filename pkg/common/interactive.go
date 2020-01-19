package common

import (
	"github.com/ktr0731/go-fuzzyfinder"
)

// InteractiveMode displays a UI that provide fuzzy finding against to the options passed
func InteractiveMode(options []string) (string, error) {
	idx, err := fuzzyfinder.Find(options, func(i int) string { return options[i] })
	if err != nil {
		return "", err
	}

	return options[idx], nil
}
