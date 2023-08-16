// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

package dashboard

import (
	"embed"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/szkiba/xk6-dashboard/assets"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

func TestNewExtension(t *testing.T) {
	t.Parallel()

	var params output.Params

	params.ConfigArgument = "port=1&host=localhost"
	params.OutputType = "dashboard"

	ext, err := New(params, embed.FS{}, embed.FS{})

	assert.NoError(t, err)
	assert.NotNil(t, ext)

	assert.Equal(t, "dashboard (localhost:1) http://localhost:1", ext.Description())

	params.ConfigArgument = "period=2"

	_, err = New(params, embed.FS{}, embed.FS{})

	assert.Error(t, err)
}

func TestExtension(t *testing.T) {
	t.Parallel()

	port := getRandomPort(t)
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

	var params output.Params

	params.Logger = logrus.StandardLogger()
	params.ConfigArgument = "period=10ms&port=" + strconv.Itoa(port)

	ext, err := New(params, embed.FS{}, embed.FS{})

	assert.NoError(t, err)
	assert.NotNil(t, ext)

	assert.NoError(t, ext.Start())

	time.Sleep(time.Millisecond)

	go func() {
		sample := testSample(t, "foo", metrics.Counter, 1)

		ext.AddMetricSamples(testSampleContainer(t, sample).toArray())
	}()

	lines := readSSE(t, 7, "http://"+addr+"/events")

	assert.NotNil(t, lines)
	assert.Equal(t, "event: snapshot", lines[2])
	assert.Equal(t, "event: cumulative", lines[6])

	dataPrefix := `data: {"foo":{`

	assert.True(t, strings.HasPrefix(lines[1], dataPrefix))
	assert.True(t, strings.HasPrefix(lines[5], dataPrefix))

	assert.NoError(t, ext.Stop())
}

func TestExtension_report(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "")

	assert.NoError(t, err)
	assert.NoError(t, file.Close())

	var params output.Params

	params.Logger = logrus.StandardLogger()
	params.ConfigArgument = "period=10ms&port=0&report=" + file.Name() + ".gz"

	ext, err := New(params, embed.FS{}, assets.DirBrief())

	assert.NoError(t, err)
	assert.NotNil(t, ext)

	assert.NoError(t, ext.Start())

	time.Sleep(time.Millisecond)

	go func() {
		sample := testSample(t, "foo", metrics.Counter, 1)

		ext.AddMetricSamples(testSampleContainer(t, sample).toArray())
	}()

	assert.NoError(t, ext.Stop())

	st, err := os.Stat(file.Name() + ".gz")

	assert.NoError(t, err)

	assert.Greater(t, st.Size(), int64(1024))

	assert.NoError(t, os.Remove(file.Name()+".gz"))
}
