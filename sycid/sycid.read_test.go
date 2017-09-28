package main

import (
	"testing"

	"github.com/jteeuwen/evdev"
)

func Test_keepLines(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1 lines", args{"eins\nzwei\ndrei\nvier\n", 1}, "eins"},
		{"2 lines", args{"eins\nzwei\ndrei\nvier\n", 2}, "eins\nzwei"},
		{"3 lines", args{"eins\nzwei\ndrei\nvier\n", 3}, "eins\nzwei\ndrei"},
		{"4 lines", args{"eins\nzwei\ndrei\nvier\n", 4}, "eins\nzwei\ndrei\nvier"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keepLines(tt.args.s, tt.args.n); got != tt.want {
				t.Errorf("keepLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readRfid(t *testing.T) {
	tDev, _ := evdev.Open("/dev/input/event17")
	// ("/dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd")

	defer tDev.Close()

	type args struct {
		dev *evdev.Device
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"readBlueChip", args{tDev}, "007"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readRfid(tt.args.dev); got != tt.want {
				t.Errorf("readRfid() = %v, want %v", got, tt.want)
			}
		})
	}
}
