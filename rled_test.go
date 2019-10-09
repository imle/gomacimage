package gomacimage

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestRledFromBytes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "1006"},
		{name: "1010"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//wantFile, err := os.OpenFile(fmt.Sprintf("test/fixtures/rleds/%s.png", tt.name), os.O_RDONLY, 0755)
			//if err != nil {
			//	t.Errorf("os.OpenFile() error = %v", err)
			//	return
			//}
			//defer wantFile.Close()
			//
			//want, err := png.Decode(wantFile)
			//if err != nil {
			//	t.Errorf("png.Decode() error = %v", err)
			//	return
			//}

			binaryData, err := ioutil.ReadFile(fmt.Sprintf("test/fixtures/rleds/%s.bin", tt.name))
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			got, err := RledFromBytes(binaryData)
			if err != nil {
				t.Errorf("RledFromBytes() error = %v", err)
				return
			}

			//err = fuzzyCompImage(got, want)
			//if err != nil {
			//	t.Errorf("fuzzyCompImage() error = %v", err)
			//}
		})
	}
}
