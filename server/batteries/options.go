/*
Copyright 2024 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package batteries

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	genericfeatures "k8s.io/apiserver/pkg/features"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

// Options holds the configuration for the batteries.
type Options struct {
	batteries BatteriesList
	Enabled   []string
}

type completedOptions struct {
	batteries BatteriesList
	Enabled   []string
}

// CompletedOptions holds the completed configuration for the batteries.
type CompletedOptions struct {
	*completedOptions
}

// AddFlags adds the flags for the admin authentication to the given FlagSet.
func (s *Options) AddFlags(fs *pflag.FlagSet) {
	if s == nil {
		return
	}

	all := sets.NewString()
	enabled := sets.NewString()
	var maxLen int
	for name := range defaultBatteries {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}
	for name, bat := range defaultBatteries {
		if bat.Groups == nil {
			all = all.Insert(fmt.Sprintf("%-*s: %s", maxLen, name, bat.Description))
		}
		all = all.Insert(fmt.Sprintf("%-*s %s [%s]", maxLen+1, name+":", bat.Description, strings.Join(bat.Groups, ", ")))
		if bat.Enabled {
			enabled.Insert(string(name))
		}
	}
	fs.StringSliceVar(&s.Enabled, "batteries", []string{}, fmt.Sprintf(
		"The batteries to enable in the generic control-plane server ('-battery' to disable, '+battery' or 'battery' to enable).\n\nPossible values:\n- %s\n\nEnabled batteries: %s",
		strings.Join(all.List(), "\n- "),
		strings.Join(enabled.List(), ", "),
	))
}

// Complete defaults fields that have not set by the consumer of this package.
func (b Options) Complete() CompletedOptions {
	// Ensure all related configurations are configured
	for _, name := range b.Enabled {
		if len(name) == 0 {
			continue
		}
		switch name[0] {
		case '-':
			if _, ok := b.batteries[Battery(name[1:])]; !ok {
				fmt.Fprintf(os.Stderr, "Warning: unknown battery %q\n", name[1:])
			}
			b.Disable(Battery(name[1:]))
		case '+':
			if _, ok := b.batteries[Battery(name[1:])]; !ok {
				fmt.Fprintf(os.Stderr, "Warning: unknown battery %q\n", name[1:])
			}
			b.Enable(Battery(name[1:]))
		default:
			if _, ok := b.batteries[Battery(name[1:])]; !ok {
				fmt.Fprintf(os.Stderr, "Warning: unknown battery %q\n", name)
			}
			b.Enable(Battery(name))
		}
	}

	ret := CompletedOptions{
		&completedOptions{
			batteries: b.batteries,
			Enabled:   b.Enabled,
		},
	}

	// If lease is disabled, we disable APIServerIdentity
	if !ret.IsEnabled(BatteryLeases) {
		utilruntime.Must(utilfeature.DefaultMutableFeatureGate.Set(fmt.Sprintf("%s=false", genericfeatures.APIServerIdentity)))
	}

	return ret
}

// Validate validates the batteries options.
func (b CompletedOptions) Validate() []error {
	var errs []error
	for _, name := range b.Enabled {
		fmt.Println(name)
		if _, ok := b.batteries[Battery(name)]; !ok {
			errs = append(errs, fmt.Errorf("invalid battery %q", name))
		}
	}
	return errs
}
