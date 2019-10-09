package gomacimage

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestCicnFromBytes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "10000"},
		{name: "10001"},
		{name: "10002"},
		{name: "10003"},
		{name: "10004"},
		{name: "10005"},
		{name: "10006"},
		{name: "10007"},
		{name: "10008"},
		{name: "10009"},
		{name: "10010"},
		{name: "10011"},
		{name: "10012"},
		{name: "10013"},
		{name: "10014"},
		{name: "10015"},
		{name: "10016"},
		{name: "10017"},
		{name: "10018"},
		{name: "10019"},
		{name: "10020"},
		{name: "10021"},
		{name: "10022"},
		{name: "10023"},
		{name: "15000"},
		{name: "15001"},
		{name: "18000"},
		{name: "18001"},
		{name: "20000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantFile, err := os.OpenFile(fmt.Sprintf("test/fixtures/cicns/%s.png", tt.name), os.O_RDONLY, 0755)
			if err != nil {
				t.Errorf("os.OpenFile() error = %v", err)
				return
			}
			defer wantFile.Close()

			want, err := png.Decode(wantFile)
			if err != nil {
				t.Errorf("png.Decode() error = %v", err)
				return
			}

			binaryData, err := ioutil.ReadFile(fmt.Sprintf("test/fixtures/cicns/%s.bin", tt.name))
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			got, err := CicnFromBytes(binaryData)
			if err != nil {
				t.Errorf("CicnFromBytes() error = %v", err)
				return
			}

			err = fuzzyCompImage(got, want)
			if err != nil {
				t.Errorf("fuzzyCompImage() error = %v", err)
			}
		})
	}
}
