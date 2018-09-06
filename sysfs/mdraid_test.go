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

	expected := []MdraidStat{
		// device, level, array state, metadata ver, total disks, chunk size, degraded disks, sync action, sync completed, mismatch count
		{"md0", "raid0", "active", "1.2", 4, 524288, 0, "", 0, 0},
		{"md1", "raid1", "clean", "1.2", 2, 0, 0, "idle", 0, 0},
		{"md10", "raid10", "clean", "1.2", 4, 524288, 0, "idle", 0, 0},
		{"md4", "raid4", "clean", "1.2", 4, 524288, 0, "idle", 0, 0},
		{"md5", "raid5", "clean", "1.2", 4, 524288, 0, "idle", 0.9920517758931722, 0},
		{"md6", "raid6", "clean", "1.2", 4, 524288, 0, "idle", 0, 0},
		{"md99", "linear", "clean", "1.2", 4, 0, 0, "", 0, 0},
	}

	if !reflect.DeepEqual(expected, stats) {
		t.Errorf("Result not correct: want %v, have %v", expected, stats)
	}
}
