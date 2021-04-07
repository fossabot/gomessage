package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func Test_initLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestLogOutputFormat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initLog()
		})
	}
}

func TestLogOutputFormat(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.Error("Helloerror")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "Helloerror", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func Test_ping(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ping(tt.args.addr)
		})
	}
}

func Test_wait(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wait()
		})
	}
}
