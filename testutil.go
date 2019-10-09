package gomacimage

import (
	"errors"
	"fmt"
	"image"
	"reflect"
)

func fuzzyCompImage(got image.Image, want image.Image) error {
	if !reflect.DeepEqual(got.Bounds(), want.Bounds()) {
		return errors.New(fmt.Sprintf("RledFromBytes() [Bounds] got = %v, want %v", got.Bounds(), want.Bounds()))
	}
	for y := 0; y < got.Bounds().Max.Y; y++ {
		for x := 0; x < got.Bounds().Max.X; x++ {
			g := got.At(x, y)
			w := want.At(x, y)
			Rg, Gg, Bg, Ag := g.RGBA()
			Rw, Gw, Bw, Aw := w.RGBA()

			if Ag == 0 && Aw == 0 {
				continue
			}

			if Rg != Rw || Bg != Bw || Gg != Gw || Ag != Aw {
				return errors.New(fmt.Sprintf("RledFromBytes() [At(%v, %v)] got = %v, want %v", x, y, g, w))
			}
		}
	}

	return nil
}
