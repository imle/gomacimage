package gopict

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestPictFromBytes(t *testing.T) {
	tests := []struct {
		name        string
		binName     string
		comapreName string
	}{
		{
			name:        "ship",
			binName:     "ship",
			comapreName: "ship",
		},
		{
			name:        "landed",
			binName:     "landed",
			comapreName: "landed",
		},
		{
			name:        "status bar",
			binName:     "statusBar",
			comapreName: "statusBar",
		},
		{
			name:        "target image",
			binName:     "targetImage",
			comapreName: "targetImage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binaryData, err := ioutil.ReadFile(fmt.Sprintf("fixtures/bins/%s.bin", tt.binName))
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			wantFile, err := os.OpenFile(fmt.Sprintf("fixtures/picts/%s.png", tt.comapreName), os.O_RDONLY, 0755)
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

			got, err := FromBytes(binaryData)
			if err != nil {
				t.Errorf("FromBytes() error = %v", err)
				return
			}

			//o, _ := os.OpenFile(fmt.Sprintf("out/%s.png", tt.name), os.O_WRONLY|os.O_CREATE, 0600)
			//defer o.Close()
			//png.Encode(o, got)

			if !reflect.DeepEqual(got.Bounds(), want.Bounds()) {
				t.Errorf("FromBytes() [Bounds] got = %v, want %v", got.Bounds(), want.Bounds())
			}

			for y := 0; y < got.Bounds().Max.Y; y++ {
				for x := 0; x < got.Bounds().Max.X; x++ {
					g := got.At(x, y)
					w := want.At(x, y)
					Rg, Gg, Bg, Ag := g.RGBA()
					Rw, Gw, Bw, Aw := w.RGBA()
					if Rg != Rw || Bg != Bw || Gg != Gw || Ag != Aw {
						t.Errorf("FromBytes() [At(%v, %v)] got = %v, want %v", x, y, g, w)
					}
				}
			}
		})
	}
}
