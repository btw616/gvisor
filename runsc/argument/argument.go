// Copyright 2020 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package argument provides mechanisms for developers to specify extra arguments
// for the boot subcommand that are not explicitly written in boot.go. This is a
// fairly restrictive pipeline only meant to allow for the specification of top-level
// runsc args that result in additional arguments being passed to the boot subcommand;
// it may be generalized further in the future.
package argument

import (
	"os"

	"gvisor.dev/gvisor/runsc/flag"
)

// Argument provides an interface for devs to add their own arguments to runsc.
type Argument interface {
	// SetFlags adds the command line flag to a flagset, such that
	// flagset.Parse() will parse the argument appropriately.
	SetFlags(f *flag.FlagSet)
	// OnCreateSandboxProcess is a hook that is evaluated near the end of
	// createSandboxProcess(). Any strings returned in extraArgs will be appended
	// to cmd.Args, and any files returned in extraFiles will be appended to
	// cmd.ExtraFiles. Implementations are responsible for updating nextFD.
	OnCreateSandboxProcess(id string, nextFD *int) (extraArgs []string, extraFiles []*os.File, err error)
	// OnBoot will be evaluated when a new kernel loader is created.
	OnBoot() error
}

// NoopArgument contains default implementations of methods that may not
// necessarily be implemented by all Argument implementations.
type NoopArgument struct {
}

// OnCreateSandboxProcess is a default implementation of
// Argument.OnCreateSandboxProcess()
func (b *NoopArgument) OnCreateSandboxProcess(id string, nextFD *int) ([]string, []*os.File, error) {
	return []string{}, []*os.File{}, nil
}

// OnBoot is a default implementation of Argument.OnBoot()
func (b *NoopArgument) OnBoot() error {
	return nil
}

// ArgSet wraps a list of Arguments, providing some convenience methods.
type ArgSet struct {
	args []Argument
}

// Register adds a an argument to the ArgSet, so that it will be set
// when SetFlags is called later.
func (a *ArgSet) Register(arg Argument) {
	a.args = append(a.args, arg)
}

// SetFlags adds an argument to the provided flagset, such that it will be
// evaluated when f.Parse() is called later.
func (a *ArgSet) SetFlags(f *flag.FlagSet) {
	for _, arg := range a.args {
		arg.SetFlags(f)
	}
}

// All returns the list of arguments wrapped by the ArgSet.
func (a *ArgSet) All() []Argument {
	return a.args
}
