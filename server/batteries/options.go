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

	"github.com/spf13/pflag"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	genericfeatures "k8s.io/apiserver/pkg/features"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

// AddFlags adds the flags for the admin authentication to the given FlagSet.
func (s *Batteries) AddFlags(fs *pflag.FlagSet) {
	if s == nil {
		return
	}

	fs.StringSliceVar(&s.BatteriesArgs, "batteries", []string{}, "The batteries to enable in the generic control-plane server.")
}

func (b Batteries) Complete() {
	// Ensure all some related configuration are configured

	for _, name := range b.BatteriesArgs {
		if _, ok := b.list[Battery(name)]; ok {
			b.Enable(Battery(name))
		}
	}

	// If lease is disabled, we disable APIServerIdentity
	if !b.IsEnabled(BatteryLeases) {
		utilruntime.Must(utilfeature.DefaultMutableFeatureGate.Set(fmt.Sprintf("%s=false", genericfeatures.APIServerIdentity)))
	}
}

func (b Batteries) Validate() []error {
	var errs []error
	for _, name := range b.BatteriesArgs {
		fmt.Println(name)
		if _, ok := b.list[Battery(name)]; !ok {
			errs = append(errs, fmt.Errorf("invalid battery %q", name))
		}
	}
	return errs
}
