package main

import (
	"net"
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

func Test_configureAmqp(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configureAmqp()
		})
	}
}

func Test_publishMessages(t *testing.T) {
	type args struct {
		messages []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publishMessages(tt.args.messages)
		})
	}
}

func Test_msgHandler(t *testing.T) {
	type args struct {
		src *net.UDPAddr
		n   int
		b   []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgHandler(tt.args.src, tt.args.n, tt.args.b)
		})
	}
}

func Test_configureUDPListener(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configureUDPListener()
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
