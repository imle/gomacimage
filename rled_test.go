package gomacimage

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestRleFromBytes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "1006"},
		{name: "1010"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wantFile, err := os.OpenFile(fmt.Sprintf("test/fixtures/rle/%s.png", tt.name), os.O_RDONLY, 0755)
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
			spriteMap := want.(*image.NRGBA)

			binaryData, err := ioutil.ReadFile(fmt.Sprintf("test/fixtures/rle/%s.bin", tt.name))
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			got, err := RleFromBytes(binaryData)
			if err != nil {
				t.Errorf("RleFromBytes() error = %v", err)
				return
			}

			writeImage(spriteMap, fmt.Sprintf("test/fixtures/rle/cmps/%v-in.png", tt.name))
			writeImage(got.Image, fmt.Sprintf("test/fixtures/rle/cmps/%v-gen.png", tt.name))

			diffGot, diffWant, errs := fuzzyCompImage(got.Image, spriteMap)
			if len(errs) != 0 {
				t.Errorf("fuzzyCompImage():\n  %v", errs)
			}

			if diffGot != nil {
				writeImage(diffGot, fmt.Sprintf("test/fixtures/rle/cmps/%v-diff-got.png", tt.name))
			}

			if diffWant != nil {
				writeImage(diffWant, fmt.Sprintf("test/fixtures/rle/cmps/%v-diff-want.png", tt.name))
			}
		})
	}
}
