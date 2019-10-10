package gomacimage

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"reflect"
)

func fuzzyCompImage(got image.Image, want image.Image) (diffGot, diffWant *image.RGBA, errs []error) {
	diffGot = image.NewRGBA(got.Bounds())
	diffWant = image.NewRGBA(want.Bounds())

	errs = make([]error, 0)

	if !reflect.DeepEqual(got.Bounds().Size(), want.Bounds().Size()) {
		errs = append(errs, errors.New(fmt.Sprintf("RleFromBytes() [Bounds] got = %v, want %v", got.Bounds(), want.Bounds())))
		return nil, nil, errs
	}
	for x := 0; x < got.Bounds().Max.X; x++ {
		for y := 0; y < got.Bounds().Max.Y; y++ {
			g := got.At(x+got.Bounds().Min.X, y+got.Bounds().Min.Y)
			w := want.At(x+want.Bounds().Min.X, y+want.Bounds().Min.Y)
			Rg, Gg, Bg, Ag := g.RGBA()
			Rw, Gw, Bw, Aw := w.RGBA()

			if Ag == 0 && Aw == 0 {
				continue
			}

			if Rg != Rw || Bg != Bw || Gg != Gw || Ag != Aw {
				diffGot.Set(x, y, g)
				diffWant.Set(x, y, w)

				errs = append(errs, errors.New(fmt.Sprintf("RleFromBytes() [At(%v, %v)] got = %v, want %v", x, y, g, w)))
			}
		}
	}

	return diffGot, diffWant, errs
}

func writeImage(img image.Image, path string) {
	gotOut, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	defer gotOut.Close()
	png.Encode(gotOut, img)
}
