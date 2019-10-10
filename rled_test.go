package gomacimage

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"sync"
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
			if err != nil && got != nil {
				t.Errorf("RleFromBytes() error = %v", err)
				return
			}

			spriteSize := got[0].Bounds().Size()
			mapSize := want.Bounds().Size()

			countAcross := mapSize.X / spriteSize.X

			writeImage(spriteMap, fmt.Sprintf("test/fixtures/rle/cmps/%v-0-sm.png", tt.name))

			jobs := make(chan job, 10)
			defer close(jobs)
			results := make(chan jobResult, 10)
			defer close(results)

			for w := 1; w <= 10; w++ {
				go worker(jobs, results)
			}

			wg := sync.WaitGroup{}

			go func() {
				for i := range got {
					result := <-results
					if len(result.errs) != 0 {
						t.Errorf("fuzzyCompImage() [sprite %03d]:\n  %v", i, result.errs)
					}
					wg.Done()
				}
			}()

			for i, v := range got {
				wg.Add(1)

				xTopLeft := i % countAcross * spriteSize.X
				yTopLeft := i / countAcross * spriteSize.Y

				rect := image.Rect(xTopLeft, yTopLeft, xTopLeft+spriteSize.X, yTopLeft+spriteSize.Y)

				jobs <- job{
					built:    v,
					subImage: spriteMap.SubImage(rect),
				}
			}

			wg.Wait()
		})
	}
}

type job struct {
	built    image.Image
	subImage image.Image
}

type jobResult struct {
	errs []error
}

func worker(jobs <-chan job, results chan<- jobResult) {
	for j := range jobs {
		_, _, errs := fuzzyCompImage(j.built, j.subImage)
		results <- jobResult{errs: errs}
	}
}
