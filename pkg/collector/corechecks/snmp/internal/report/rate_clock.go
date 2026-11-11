// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package report

import (
	"math"
	"time"
)

// rateClockMaxDrift is the largest difference, in seconds, that the clock accepts between the
// time from sysUpTime and the agent clock before it uses the agent clock again.
const rateClockMaxDrift = 2.0

// RateClock gives the timestamp of the rate samples of one device.
//
// The check reports the values only after the whole fetch, and the fetch duration changes from one
// run to the next. A rate must use the time between two reads of the device, so the clock moves
// with the device sysUpTime. It uses the start time of the fetch when sysUpTime cannot be used:
// on the first run, when the value is missing, when sysUpTime decreases (reset or wrap), and when
// the time from sysUpTime drifts from the agent clock by more than rateClockMaxDrift.
type RateClock struct {
	lastTimestamp float64 // seconds since epoch
	lastUptime    float64 // hundredths of a second
	hasLast       bool
}

// NewRateClock returns a new RateClock
func NewRateClock() *RateClock {
	return &RateClock{}
}

// Next returns the timestamp, in seconds since epoch, of the samples of a fetch that started at
// fetchTime. uptime is the sysUpTime value of that fetch, in hundredths of a second.
func (c *RateClock) Next(fetchTime time.Time, uptime float64, hasUptime bool) float64 {
	fetchTs := float64(fetchTime.UnixNano()) / float64(time.Second)
	ts := fetchTs
	if c.hasLast && hasUptime && uptime > c.lastUptime {
		uptimeTs := c.lastTimestamp + (uptime-c.lastUptime)/100
		if math.Abs(uptimeTs-fetchTs) <= rateClockMaxDrift {
			ts = uptimeTs
		}
	}
	c.lastTimestamp, c.lastUptime, c.hasLast = ts, uptime, hasUptime
	return ts
}
