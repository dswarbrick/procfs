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
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// MdraidStat holds info parsed from various files in the /sys/block/md*/md directory.
type MdraidStat struct {
	Device          string // Kernel device name of array.
	Level           string // mdraid level.
	ArrayState      string // State of the array.
	MetadataVersion string // mdraid metadata version.
	TotalDisks      uint64 // Total number of disks in the array.
	ChunkSize       uint64 // Chunk size.

	// The following items are only valid for raid1, 4, 5, 6 and 10.
	DegradedDisks uint64  // Number of degraded disks in the array.
	SyncAction    string  // Current sync action.
	SyncCompleted float64 // Fraction (0.0 - 1.0) representing the completion status of current sync operation.
	MismatchCount uint64  // Number of mdraid mismatches.
}

func (fs FS) NewMdraidStat() ([]MdraidStat, error) {
	matches, err := filepath.Glob(fs.Path("block/md*/md"))
	if err != nil {
		return nil, err
	}

	stats := make([]MdraidStat, 0, len(matches))

	for _, m := range matches {
		md := MdraidStat{Device: filepath.Base(filepath.Dir(m))}
		path := fs.Path("block", md.Device, "md")

		if val, err := sysReadFileString(path + "/level"); err == nil {
			md.Level = val
		} else {
			return nil, err
		}

		if val, err := sysReadFileString(path + "/array_state"); err == nil {
			md.ArrayState = val
		} else {
			return nil, err
		}

		if val, err := sysReadFileUint64(path + "/chunk_size"); err == nil {
			md.ChunkSize = val
		} else {
			return nil, err
		}

		if val, err := sysReadFileString(path + "/metadata_version"); err == nil {
			md.MetadataVersion = val
		} else {
			return nil, err
		}

		if val, err := sysReadFileUint64(path + "/raid_disks"); err == nil {
			md.TotalDisks = val
		} else {
			return nil, err
		}

		switch md.Level {
		case "raid1", "raid4", "raid5", "raid6", "raid10":
			if val, err := sysReadFileUint64(path + "/degraded"); err == nil {
				md.DegradedDisks = val
			} else {
				return nil, err
			}

			if val, err := sysReadFileUint64(path + "/mismatch_cnt"); err == nil {
				md.MismatchCount = val
			} else {
				return nil, err
			}

			if val, err := sysReadFileString(path + "/sync_action"); err == nil {
				md.SyncAction = val
			} else {
				return nil, err
			}

			if val, err := sysReadFileString(path + "/sync_completed"); err == nil {
				if val != "none" {
					var a, b uint64

					if _, err := fmt.Sscanf(val, "%d / %d", &a, &b); err == nil {
						md.SyncCompleted = float64(a) / float64(b)
					} else {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		}

		stats = append(stats, md)
	}

	return stats, nil
}

func sysReadFileString(file string) (string, error) {
	fileContents, err := sysReadFile(file)
	if err != nil {
		return "", err
	}

	s := strings.TrimSpace(string(fileContents))
	return s, nil
}

func sysReadFileUint64(file string) (uint64, error) {
	fileContents, err := sysReadFile(file)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSpace(string(fileContents)), 10, 64)
}
