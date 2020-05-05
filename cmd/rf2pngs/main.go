package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/imle/resourcefork"

	"github.com/imle/gomacimage"
)

func main() {
	dir := "test/fixtures/rle2"
	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	rf, err := resourcefork.ReadResourceForkFromPath("./assets/Nova Files")
	if err != nil {
		log.Fatal(err)
	}

	for id, res := range rf.Resources["rlëD"] {
		func() {
			fmt.Println(id)
			want, err := gomacimage.RleFromBytes(res.Data)
			if err != nil {
				log.Fatalf("RleFromBytes() error = %v", err)
			}

			o, _ := os.OpenFile(fmt.Sprintf("%s/%d.png", dir, id), os.O_WRONLY|os.O_CREATE, 0600)
			defer o.Close()
			_ = png.Encode(o, want.Image)
		}()
	}
}
