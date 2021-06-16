// Copyright (c) 2017 Sweetbridge Inc.
// Copyright (c) 2018 Robert Zaremba
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

package setup

import (
	"fmt"
	"os"
	"strings"

	"github.com/robert-zaremba/flag"
)

// FlagFail - exits the main process and displays usage information
func FlagFail(err error) {
	logger.Error("!! Wrong CMD parameters !! Run with `-h` parameter to output a usage information", err)

	os.Exit(1)
}

// FlagValidate runs flag checkers
func FlagValidate(positionalArgs string, checkers ...Checker) {
	expected := strings.Fields(positionalArgs)
	if len(expected) != flag.NArg() {
		FlagFail(fmt.Errorf("missing required positional arguments: %s. Number of args Provided: %d, expected: %d",
			expected, flag.NArg(), len(expected)))
	}
	for _, c := range checkers {
		if err := c.Check(); err != nil {
			FlagFail(err)
		}
	}
}

// FlagCheckMany is a helper function to check many flag components
func FlagCheckMany(checkers ...Checker) error {
	for _, c := range checkers {
		if err := c.Check(); err != nil {
			return err
		}
	}
	return nil
}

// FlagSimpleInit provides a common functionality to setup the command line flags.
// `positionalArgs` documents expected positional argumentes, eg `"arg1 arg2 arg3"`.
// `rollbarKey` is a pointer, because it can be a flag, which is gonig to be initialized
//    in this function.
func FlagSimpleInit(name, positionalArgs string, flags ...Checker) {
	flag.Setup(GitVersion, positionalArgs)
	FlagValidate(positionalArgs, flags...)
	var rbKey string
	if rollbarFlag != nil {
		rbKey = *rollbarFlag
	}
	MustLogger(name, rbKey)
}

// Checker is an interface for type which has a Check function
type Checker interface {
	Check() error
}
