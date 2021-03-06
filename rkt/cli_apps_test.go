// Copyright 2014 CoreOS, Inc.
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

package main

import (
	"flag"
	"reflect"
	"strings"
	"testing"
)

func TestParseAppArgs(t *testing.T) {
	flags := flag.NewFlagSet("test", flag.ExitOnError)
	tests := []struct {
		in     string
		images []string
		args   [][]string
		werr   bool
	}{
		{
			"example.com/foo example.com/bar -- --help --- example.com/baz -- --verbose",
			[]string{"example.com/foo", "example.com/bar", "example.com/baz"},
			[][]string{
				nil,
				[]string{"--help"},
				[]string{"--verbose"},
			},
			false,
		},
		{
			"example.com/foo --- example.com/bar --- example.com/baz ---",
			[]string{"example.com/foo", "example.com/bar", "example.com/baz"},
			[][]string{
				nil,
				nil,
				nil,
			},
			false,
		},
		{
			"example.com/foo example.com/bar example.com/baz",
			[]string{"example.com/foo", "example.com/bar", "example.com/baz"},
			[][]string{
				nil,
				nil,
				nil,
			},
			false,
		},
	}

	for i, tt := range tests {
		rktApps.Reset()
		err := parseApps(&rktApps, strings.Split(tt.in, " "), flags, true)
		ga := rktApps.GetArgs()
		gi := rktApps.GetImages()
		if gerr := (err != nil); gerr != tt.werr {
			t.Errorf("#%d: err==%v, want errstate %t", i, err, tt.werr)
		}
		if !reflect.DeepEqual(ga, tt.args) {
			t.Errorf("#%d: got args %v, want args %v", i, ga, tt.args)
		}
		if !reflect.DeepEqual(gi, tt.images) {
			t.Errorf("#%d: got images %v, want images %v", i, gi, tt.images)
		}
	}

}
