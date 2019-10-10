package main

import (
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/imle/gomacimage"
)

func main() {
	binaryData, err := ioutil.ReadFile("test/fixtures/rleds/1006.bin")
	if err != nil {
		log.Fatalf("ioutil.ReadFile() error = %v", err)
	}

	got, err := gomacimage.StitchedRledFromBytes(binaryData, 6)
	if err != nil && got != nil {
		log.Fatalf("RledFromBytes() error = %v", err)
	}

	o, _ := os.OpenFile("/Users/ski/go/src/gomacimage/test/fixtures/rleds/1006.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer o.Close()
	png.Encode(o, got)
}
