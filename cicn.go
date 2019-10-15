package gomacimage

import (
	"errors"
	"fmt"
	"image"
	"image/color"
)

func CicnFromBytes(b []byte) (img image.Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			img = nil

			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New(fmt.Sprintf("unknown panic: %v", r))
			}
		}
	}()

	parser := dataStructureParse{
		d:   NewBigEndianDataView(b),
		pos: 0,
	}

	pixelMap := parser.parsePixMap()
	maskBitMap := parser.parseBitMap()
	iconBitMap := parser.parseBitMap()
	_ = parser.readDWord() // Should always be 0x00000000

	maskBitMapImageDataLength := maskBitMap.rowBytes * maskBitMap.bounds.height
	maskBitMapImageData := parser.readDataUint8(int(maskBitMapImageDataLength))

	iconBitMapImageDataLength := iconBitMap.rowBytes * iconBitMap.bounds.height
	_ = parser.readDataUint8(int(iconBitMapImageDataLength))

	colorTable := parser.parseColorTable()

	pixelMapImageDataLength := pixelMap.rowBytes * pixelMap.bounds.height
	pixelMapImageData := parser.readDataUint8(int(pixelMapImageDataLength))

	rect := pixelMap.bounds
	imgRGBA := image.NewNRGBA(image.Rect(int(rect.x), int(rect.y), int(rect.width), int(rect.height)))
	for x := 0; x < int(rect.width); x++ {
		for y := 0; y < int(rect.height); y++ {
			idx := uint32(y)*uint32(pixelMap.rowBytes&0x3FFF)*8/uint32(pixelMap.pixelSize) + uint32(x)

			var col uint16
			switch pixelMap.pixelSize {
			case 1:
				col = uint16(pixelMapImageData[idx/8])
				col &= 0x80 >> (idx % 8)
				col >>= 7 - (idx % 8)
			case 2:
				col = uint16(pixelMapImageData[idx/4])
				switch idx % 4 {
				case 0:
					col >>= 2
					fallthrough
				case 1:
					col >>= 2
					fallthrough
				case 2:
					col >>= 2
					fallthrough
				case 3:
					col &= 0x03
				}
			case 4:
				col = uint16(pixelMapImageData[idx/2])
				if idx%2 == 1 {
					col &= 0x0F
				} else {
					col = (col & 0xF0) >> 4
				}
			case 8:
				col = uint16(pixelMapImageData[idx])
			default:
				return nil, errors.New(fmt.Sprintf("unhandled pixel size: %v", pixelMap.pixelSize))
			}

			for i := uint16(0); i < colorTable.size; i++ {
				if colorTable.data[i].value == col {
					imgRGBA.Set(x, y, color.RGBA{
						R: uint8(colorTable.data[i].r),
						G: uint8(colorTable.data[i].g),
						B: uint8(colorTable.data[i].b),
						A: ((maskBitMapImageData[idx/8] >> uint8(7-idx%8)) & 0x1) * 255,
					})
					break
				}
			}
		}
	}

	return imgRGBA, nil
}
