//  Copyright 2020 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteralLoggable_GetValueAsInt64Slice(t *testing.T) {
	loggable := literalLoggable{
		int64s: map[string][]int64{
			"gb": {1, 2, 3},
		},
	}

	assert.Equal(t, []int64{1, 2, 3}, loggable.GetValueAsInt64Slice("gb"))
	assert.Empty(t, loggable.GetValueAsInt64Slice("not-there"))
}

func TestLiteralLoggable_GetValue(t *testing.T) {
	loggable := literalLoggable{
		strings: map[string]string{"hello": "world"},
	}

	assert.Equal(t, "world", loggable.GetValue("hello"))
	assert.Empty(t, loggable.GetValue("not-there"))
}

func TestLiteralLoggable_ReadSerialPortLogs(t *testing.T) {
	loggable := literalLoggable{
		traceLogs: []string{"log-a", "log-b"},
	}

	assert.Equal(t, []string{"log-a", "log-b"}, loggable.ReadSerialPortLogs())
}

func TestSingleImageImportLoggableBuilder(t *testing.T) {
	format := "vmdk"
	sourceGb := int64(12)
	targetGb := int64(100)
	traceLogs1 := []string{"log-a", "log-b"}
	traceLogs2 := []string{"log-c", "log-d"}
	inflationTypeValue := "qemu"
	inflationTimeValue := int64(10000)
	shadowInflationTimeValue := int64(5000)
	matchResultValue := "true"
	bootFSValue := "btrfs"
	for _, isUEFICompatibleImageValue := range []bool{true, false} {
		for _, isUEFIDetectedValue := range []bool{true, false} {
			for _, biosBootableValue := range []bool{true, false} {
				expected := literalLoggable{
					strings: map[string]string{
						importFileFormat:      format,
						inflationType:         inflationTypeValue,
						shadowDiskMatchResult: matchResultValue,
						rootFS:                bootFSValue,
					},
					int64s: map[string][]int64{
						sourceSizeGb:        {sourceGb},
						targetSizeGb:        {targetGb},
						inflationTime:       {inflationTimeValue},
						shadowInflationTime: {shadowInflationTimeValue},
					},
					bools: map[string]bool{
						isUEFICompatibleImage: isUEFICompatibleImageValue,
						isUEFIDetected:        isUEFIDetectedValue,
						uefiBootable:          isUEFIDetectedValue,
						biosBootable:          biosBootableValue,
					},
					traceLogs: append(traceLogs1, traceLogs2...),
				}
				assert.Equal(t, expected, NewSingleImageImportLoggableBuilder().
					SetDiskAttributes(format, sourceGb, targetGb).
					SetUEFIMetrics(isUEFICompatibleImageValue, isUEFIDetectedValue, biosBootableValue, bootFSValue).
					SetInflationAttributes(matchResultValue, inflationTypeValue, inflationTimeValue, shadowInflationTimeValue).
					AppendTraceLogs(traceLogs1).
					AppendTraceLogs(traceLogs2).
					Build())
			}
		}
	}
}

func TestOvfExportLoggableBuilder(t *testing.T) {
	sourceGb := []int64{12, 25}
	targetGb := []int64{100, 300}
	traceLogs1 := []string{"log-a", "log-b"}
	traceLogs2 := []string{"log-c", "log-d"}

	expected := literalLoggable{
		strings: map[string]string{},
		int64s: map[string][]int64{
			sourceSizeGb: {12, 25},
			targetSizeGb: {100, 300},
		},
		bools:     map[string]bool{},
		traceLogs: append(traceLogs1, traceLogs2...),
	}
	assert.Equal(t, expected, NewOvfExportLoggableBuilder().SetDiskSizes(sourceGb, targetGb).
		AppendTraceLogs(traceLogs1).
		AppendTraceLogs(traceLogs2).
		Build())
}
