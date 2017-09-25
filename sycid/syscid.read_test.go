package main

import (
	"testing"

	"github.com/dddpaul/golang-evdev/evdev"
)

func Test_readRfid(t *testing.T) {
	type args struct {
		fd *evdev.InputDevice
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRfid(tt.args.fd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readRfid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readRfid() = %v, want %v", got, tt.want)
			}
		})
	}
}
