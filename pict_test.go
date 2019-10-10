package gomacimage

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestPictFromBytes(t *testing.T) {
	tests := []struct {
		name        string
		binName     string
		compareName string
	}{
		{
			name:        "ship",
			binName:     "ship",
			compareName: "ship",
		},
		{
			name:        "landed",
			binName:     "landed",
			compareName: "landed",
		},
		{
			name:        "status bar",
			binName:     "statusBar",
			compareName: "statusBar",
		},
		{
			name:        "target image",
			binName:     "targetImage",
			compareName: "targetImage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantFile, err := os.OpenFile(fmt.Sprintf("test/fixtures/picts/%s.png", tt.compareName), os.O_RDONLY, 0755)
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

			binaryData, err := ioutil.ReadFile(fmt.Sprintf("test/fixtures/picts/%s.bin", tt.binName))
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			got, err := PictFromBytes(binaryData)
			if err != nil {
				t.Errorf("PictFromBytes() error = %v", err)
				return
			}

			_, _, errs := fuzzyCompImage(got, want)
			if len(errs) != 0 {
				for _, err := range errs {
					t.Errorf("fuzzyCompImage() error = %v", err)
				}
			}
		})
	}
}
