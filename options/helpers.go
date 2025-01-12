package options

import (
	"fmt"
	"slices"
	"strings"

	"github.com/akamensky/argparse"
)

func selectorList(required bool, allowed []string, help string) *argparse.Options {
	return &argparse.Options{
		Required: required,
		Validate: func(args []string) error {
			seen := map[string]*struct{}{}

			for _, choice := range strings.Split(args[0], ",") {
				if !slices.Contains(allowed, choice) {
					return fmt.Errorf("invalid '%s'", choice)
				}
				if _, ok := seen[choice]; ok {
					return fmt.Errorf("double '%s'", choice)
				}
				seen[choice] = nil
			}

			return nil
		},
		Help: fmt.Sprintf("%s. {%s}", help, strings.Join(allowed, ",")),
	}
}
