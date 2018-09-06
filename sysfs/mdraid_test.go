// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysfs

import (
	"reflect"
	"testing"
)

func TestNewMdraidStat(t *testing.T) {
	fs, err := NewFS("fixtures")
	if err != nil {
		t.Fatal(err)
	}

	stats, err := fs.NewMdraidStat()
	if err != nil {
		t.Fatal(err)
	}

	expected := []MDStat{
		{"md0", "active", 4, 0, "raid0", "1.2", 0, "", 524288, 0},
		{"md1", "clean", 2, 0, "raid1", "1.2", 0, "idle", 0, 0},
		{"md10", "clean", 4, 0, "raid10", "1.2", 0, "idle", 524288, 0},
		{"md4", "clean", 4, 0, "raid4", "1.2", 0, "idle", 524288, 0},
		{"md5", "clean", 4, 0, "raid5", "1.2", 0, "idle", 524288, 0.9920517758931722},
		{"md6", "clean", 4, 0, "raid6", "1.2", 0, "idle", 524288, 0},
		{"md99", "clean", 4, 0, "linear", "1.2", 0, "", 0, 0},
	}

	if !reflect.DeepEqual(expected, stats) {
		t.Errorf("Result not correct: want %v, have %v", expected, stats)
	}
}
