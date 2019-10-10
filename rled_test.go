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
			wantFile, err := os.OpenFile(fmt.Sprintf("test/fixtures/rleds/%s.png", tt.name), os.O_RDONLY, 0755)
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

			binaryData, err := ioutil.ReadFile(fmt.Sprintf("test/fixtures/rleds/%s.bin", tt.name))
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			got, err := RledFromBytes(binaryData)
			if err != nil && got != nil {
				t.Errorf("RledFromBytes() error = %v", err)
				return
			}

			spriteSize := got[0].Bounds().Size()
			mapSize := want.Bounds().Size()

			countAcross := mapSize.X / spriteSize.X

			writeImage(spriteMap, fmt.Sprintf("test/fixtures/rleds/cmps/%v-0-sm.png", tt.name))
			for i, v := range got {
				xTopLeft := i % countAcross * spriteSize.X
				yTopLeft := i / countAcross * spriteSize.Y

				rect := image.Rect(xTopLeft, yTopLeft, xTopLeft+spriteSize.X, yTopLeft+spriteSize.Y)

				subImage := spriteMap.SubImage(rect)
				diffGot, diffWant, errs := fuzzyCompImage(v, subImage)
				if len(errs) != 0 {
					t.Errorf("fuzzyCompImage() [sprite %03d]:\n  %v", i, errs)
				}

				if diffGot != diffWant { // nils
					func() {
						writeImage(subImage, fmt.Sprintf("test/fixtures/rleds/cmps/%v-%03d-s.png", tt.name, i))
						writeImage(diffGot, fmt.Sprintf("test/fixtures/rleds/cmps/%v-%03d-g.png", tt.name, i))
						writeImage(diffWant, fmt.Sprintf("test/fixtures/rleds/cmps/%v-%03d-w.png", tt.name, i))
					}()
				}
			}
		})
	}
}
