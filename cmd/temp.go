package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/imle/gomacimage"
)

func main() {
	binaryData, err := ioutil.ReadFile("test/fixtures/rle/1010.bin")
	if err != nil {
		log.Fatalf("ioutil.ReadFile() error = %v", err)
	}

	want, err := gomacimage.RleFromBytes(binaryData)
	if err != nil {
		log.Fatalf("RleFromBytes() error = %v", err)
	}

	o, _ := os.OpenFile("test/fixtures/rle/1010.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer o.Close()
	png.Encode(o, want.Image)
}
