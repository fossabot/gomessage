package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/streadway/amqp"
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

func Test_processDeliveries(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processDeliveries()
		})
	}
}

func Test_decodeMessageBatch(t *testing.T) {
	type args struct {
		d []amqp.Delivery
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeMessageBatch(tt.args.d)
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
