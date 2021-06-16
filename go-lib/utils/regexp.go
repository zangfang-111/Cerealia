// Copyright (c) 2017 Sweetbridge Inc.
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

package utils

import (
	"regexp"
)

// static variables
var (
	ReAlpha    = regexp.MustCompile(`^[:alpha:]*$`)
	ReID       = regexp.MustCompile(`^[0-9A-Za-z.-]+$`)
	ReArangoID = regexp.MustCompile(`^[0-9]+$`)
	ReDBID     = regexp.MustCompile(`^[0-9A-Za-z-]+$`)
	ReAlphanum = regexp.MustCompile(`^[:alnum:]*$`)
	ReNumbers  = regexp.MustCompile(`^[0-9]*$`)
	ReEmail    = regexp.MustCompile(`^(\w[-.%+\w]*\w@\w[-.\w]*\w\.\w{2,})$`)
	RePhone    = regexp.MustCompile(`^[+]?[(]?[0-9]{1,4}[)]?[-\s0-9]*$`)
	ReHex      = regexp.MustCompile(`^[A-Fa-f0-9]+$`)
	ReSpace    = regexp.MustCompile(`\s+`)
)
