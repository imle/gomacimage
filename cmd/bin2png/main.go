package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/imle/gomacimage"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("need an id as second arg")
	}
	strId := os.Args[1]
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Fatal(err)
	}

	binaryData, err := ioutil.ReadFile(fmt.Sprintf("test/fixtures/rle/%d.bin", id))
	if err != nil {
		log.Fatalf("ioutil.ReadFile() error = %v", err)
	}

	want, err := gomacimage.RleFromBytes(binaryData)
	if err != nil {
		log.Fatalf("RleFromBytes() error = %v", err)
	}

	o, _ := os.OpenFile(fmt.Sprintf("test/fixtures/rle/%d.png", id), os.O_WRONLY|os.O_CREATE, 0600)
	defer o.Close()
	png.Encode(o, want.Image)
}
