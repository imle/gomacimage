package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/imle/resourcefork"

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

	rf, err := resourcefork.ReadResourceForkFromPath("./assets/Nova Files")
	if err != nil {
		log.Fatal(err)
	}

	res := rf.Resources["rlëD"][uint16(id)]

	want, err := gomacimage.RleFromBytes(res.Data)
	if err != nil {
		log.Fatalf("RleFromBytes() error = %v", err)
	}

	o, _ := os.OpenFile(fmt.Sprintf("test/fixtures/rle/%d.png", id), os.O_WRONLY|os.O_CREATE, 0600)
	defer o.Close()
	png.Encode(o, want.Image)
}
