// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateClock(t *testing.T) {
	start := time.Unix(1000, 0)
	at := func(seconds float64) time.Time {
		return start.Add(time.Duration(seconds * float64(time.Second)))
	}

	tests := []struct {
		name      string
		fetchTime time.Time
		uptime    float64
		hasUptime bool
		expected  float64
	}{
		{"first run uses the fetch time", at(0), 100000, true, 1000},
		{"sysUpTime moves the clock, not the fetch time", at(15.8), 101500, true, 1015},
		{"a late fetch keeps the sysUpTime interval", at(31.5), 103000, true, 1030},
		{"missing sysUpTime uses the fetch time", at(45), 0, false, 1045},
		{"the run after a missing sysUpTime uses the fetch time", at(60.3), 106000, true, 1060.3},
		{"sysUpTime is used again", at(75), 107500, true, 1075.3},
		{"a sysUpTime reset uses the fetch time", at(90), 200, true, 1090},
		{"sysUpTime is used again after a reset", at(105.4), 1700, true, 1105},
		{"a drift above the limit uses the fetch time", at(120), 6700, true, 1120},
	}

	clock := NewRateClock()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, clock.Next(tt.fetchTime, tt.uptime, tt.hasUptime), 1e-6)
		})
	}
}
