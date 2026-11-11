// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package corechecks

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DataDog/datadog-agent/pkg/aggregator/sender"
)

func TestCloneTags(t *testing.T) {
	t.Run("tags", func(t *testing.T) {
		tags := []string{"tag1", "tag2"}
		cloned := cloneTags(tags)
		tags[0] = "uh-oh"
		tags[1] = "oh no"
		require.Equal(t, []string{"tag1", "tag2"}, cloned)
	})

	t.Run("nil", func(t *testing.T) {
		cloned := cloneTags(nil)
		require.Nil(t, cloned)
	})
}

type tagRecordingSender struct {
	sender.Sender
	tags []string
}

func (s *tagRecordingSender) RateWithTimestamp(_ string, _ float64, _ string, tags []string, _ float64) error {
	s.tags = tags
	return nil
}

func TestSafeSenderRateWithTimestampClonesTags(t *testing.T) {
	recorder := &tagRecordingSender{}
	ss := newSafeSender(recorder)
	tags := []string{"a:1", "b:2"}

	require.NoError(t, ss.RateWithTimestamp("metric", 1, "", tags, 1))
	tags[0] = "changed"

	require.Equal(t, []string{"a:1", "b:2"}, recorder.tags)
}
