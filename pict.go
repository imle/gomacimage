//
// MIT License
//
// Copyright (c) 2016 Tom Hancocks, 2018 Matthew Soulanille
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// (1) Adapted from https://github.com/dmaulikr/OpenNova/blob/master/ResourceKit/ResourceFork/Parsers/RKPictureResourceParser.m

// (2) Also see http://mirrors.apple2.org.za/apple.cabi.net/Graphics/PICT.and_QT.INFO/PICT.file.format.TI.txt

// (3) Also see https://github.com/mattsoulanille/NovaParse/blob/master/src/resourceParsers/PICTParse.ts

package gomacimage

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
)

type PictOpCode uint16

const (
	PictOpCodeNop            PictOpCode = 0x0000
	PictOpCodeClipRegion                = 0x0001
	PictOpCodeDirectBitsRect            = 0x009A
	PictOpCodeEof                       = 0x00FF
	PictOpCodeDefHiLite                 = 0x001E
	PictOpCodeLongComment               = 0x00A1
	PictOpCodeExtHeader                 = 0x0C00
)

func PictFromBytes(b []byte) (img image.Image, err error) {
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

	// The first word appears to be unused so skip it.
	parser.pos += WordSize

	// The first part of the PICT is the frame.
	frame := parser.readQDRect()

	// The next 4 bytes are the version of the PICT. We're only interested in version 2.
	version := parser.readDWord()
	if version != 0x001102ff {
		return nil, errors.New("PICT resource is not version 2")
	}

	// Ensure we have an extended header here.
	opCode := parser.readOpCode()
	if opCode != PictOpCodeExtHeader {
		return nil, errors.New("expected an extended header in PICT resource")
	}

	// The next value is the header version. PICT version 2 has two variants that need to
	// be handled for EV Nova. Annoyingly it seems to use both in its data files. Not sure
	// how that happened?
	headerVersion := parser.readDWord()
	if (headerVersion & 0xFFFF0000) != 0xFFFE0000 { // Standard Header Version
		// Determine the image resolution.
		y2 := parser.readFixedPoint()
		x2 := parser.readFixedPoint()
		w2 := parser.readFixedPoint()
		h2 := parser.readFixedPoint()

		parser.xRatio = uint16(uint32(frame.x2-frame.x1) / (w2 - x2))
		parser.yRatio = uint16(uint32(frame.y2-frame.y1) / (h2 - y2))
	} else { // Extended Header Version
		parser.pos += 2 * (2 * WordSize) // 2 * uint32
		var rect = parser.readQDRect()
		parser.xRatio = (frame.x2 - frame.x1) / (rect.x2 - rect.x1)
		parser.yRatio = (frame.y2 - frame.y1) / (rect.y2 - rect.y1)
	}

	// Verify ratio is valid
	if parser.xRatio <= 0 || parser.yRatio <= 0 {
		return nil, errors.New(fmt.Sprintf("got an invalid ratio: [%v, %x]", parser.xRatio, parser.yRatio))
	}

	var op PictOpCode

	for parser.pos < len(b) {
		op = PictOpCode(parser.readOpCode())

		switch op {
		case PictOpCodeClipRegion:
			parser.readRegionWithRect()
		case PictOpCodeDirectBitsRect:
			img, err = parser.parseDirectBitsRect()
			if err != nil {
				return nil, err
			}
		case PictOpCodeLongComment:
			parser.parseLongComment()
		case PictOpCodeEof:
			return img, nil
		case PictOpCodeNop:
		case PictOpCodeExtHeader:
		case PictOpCodeDefHiLite:
		default:
			return nil, errors.New(fmt.Sprintf("encountered an unhandled opcode: [%04x]", op))
		}
	}

	return img, nil
}

func (p *dataStructureParse) readRegionWithRect() regionRect {
	var size = p.readWord()
	var regionRect = regionRect{
		x:      p.readWord() / p.xRatio,
		y:      p.readWord() / p.yRatio,
		width:  p.readWord() / p.xRatio,
		height: p.readWord() / p.yRatio,
	}
	regionRect.width -= regionRect.x
	regionRect.height -= regionRect.y
	var points = (size - 10) / 4
	p.pos += int(2 * 2 * points)
	return regionRect
}

func (p *dataStructureParse) parseLongComment() {
	var _ = p.readWord() // kind
	var length = p.readWord()
	p.pos += int(length)
}

func (p *dataStructureParse) packBitsDecode(valueSize int, data *DataView) ([]uint8, error) {
	// valueSize is in bytes, byteLength is how many bytes to read
	var result []uint8
	var pos = 0
	var length = data.GetLength()
	if valueSize > 4 {
		return nil, errors.New(fmt.Sprintf("valueSize too large. Must be <= 4 but got %v", valueSize))
	}

	var run int
	for pos < length {
		var count = data.GetUint8(pos)
		pos++

		//fmt.Printf("count: %v\n", count)

		if count < 128 {
			run = int(1+count) * valueSize
			for i := 0; i < run; i++ {
				result = append(result, data.GetUint8(pos+i))
			}
			pos += run
		} else {
			// Expand the repeat compression
			run = 256 - int(count)
			var val []uint8
			for i := 0; i < valueSize; i++ {
				val = append(val, data.GetUint8(pos+i))
			}
			pos += valueSize
			for i := 0; i <= run; i++ {
				result = append(result, val...)
			}
		}
	}

	return result, nil
}

func (p *dataStructureParse) parseDirectBitsRect() (image.Image, error) {
	px := p.parsePixMap()
	sourceRect := p.readWHRect()
	destinationRect := p.readWHRect()

	// The next 2 bytes represent the "mode" for the direct bits packing. However
	// this doesn't seem to be required with the images included in EV Nova.
	p.pos += 2

	var (
		raw          []uint8
		pxShortArray []uint16
		pxArray      []uint32
	)

	if px.packType == 3 {
		raw = make([]uint8, px.rowBytes)
	} else if px.packType == 4 {
		raw = make([]uint8, int(math.Floor(float64(int32(px.cmpCount)*int32(px.rowBytes))/4.0)))
	} else {
		return nil, errors.New(fmt.Sprintf("unsupported pack type: %v", px.packType))
	}

	pxShortArray = make([]uint16, int32(sourceRect.height)*(int32(px.rowBytes)+1))
	pxArray = make([]uint32, int(math.Floor(float64(int32(sourceRect.height)*(int32(px.rowBytes)+3))/4.0)))

	var (
		pxBufOffset      = uint32(0)
		packedBytesCount = uint16(0)
	)

	var err error
	for scanline := uint32(0); scanline < uint32(sourceRect.height); scanline++ {
		// Narrow pictures don't use the pack bits compression. Not certain what the deciding factor
		// for such a thing is, but low numbers of rowBytes seem to be the cause. Setting this to the
		// highest value found that doesn't have compression
		if px.rowBytes < 8 { // No PackBits Compression
			// gets px.rowBytes number of bytes from d
			// Then, puts sourceRect.width * 2 of them in 'raw'
			var data = p.readDataUint8(int(px.rowBytes))
			raw = data[0 : sourceRect.width*2]
		} else { // Pack Bits Compression
			if px.rowBytes > 250 {
				packedBytesCount = p.readWord()
			} else {
				packedBytesCount = uint16(p.readByte())
			}

			var encodedScanLine = p.readData(int(packedBytesCount))
			var decodedScanLine []uint8
			if px.packType == 3 {
				decodedScanLine, err = p.packBitsDecode(2, encodedScanLine)
				if err != nil {
					return nil, err
				}
			} else {
				decodedScanLine, err = p.packBitsDecode(1, encodedScanLine)
				if err != nil {
					return nil, err
				}
			}
			raw = decodedScanLine[0 : sourceRect.width*2]
		}

		if px.packType == 3 {
			// Store the decoded pixel data.
			for i := uint32(0); i < uint32(sourceRect.width); i++ {
				pxShortArray[pxBufOffset+i] = ((0xFF & uint16(raw[2*i])) << 8) | (0xFF & uint16(raw[2*i+1]))
			}
		} else {
			if px.cmpCount == 3 {
				// RGB Data
				for i := uint32(0); i < uint32(sourceRect.width); i++ {
					a := uint32(0xFF000000)
					r := (uint32(raw[i]) & 0xFF) << 16
					g := (uint32(raw[uint32(px.bounds.width)+i]) & 0xFF) << 8

					// TODO: Determine why neither of the other solutions require this.
					var b uint32
					if 2*uint32(px.bounds.width)+i >= uint32(len(raw)) {
						b = 0
					} else {
						b = uint32(raw[2*uint32(px.bounds.width)+i]) & 0xFF
					}

					pxArray[pxBufOffset+i] = a | r | g | b
				}
			} else {
				// ARGB Data
				for i := uint32(0); i < uint32(sourceRect.width); i++ {
					pxArray[pxBufOffset+i] = (uint32(raw[i])&0xFF)<<24 | (uint32(raw[uint32(px.bounds.width)+i])&0xFF)<<16 | (uint32(raw[2*uint32(px.bounds.width)+i])&0xFF)<<8 | (uint32(raw[3*uint32(px.bounds.width)+i]) & 0xFF)
				}
			}
		}

		pxBufOffset += uint32(sourceRect.width)
	}

	// Finally we need to unpack all of the pixel data. This is due to the pixels being
	// stored in an RGB 555 format. CoreGraphics does not expose a way of cleanly/publically
	// parsing this type of encoding so we need to convert it to a more modern
	// representation, such as RGBA 8888
	var (
		sourceLength = uint32(destinationRect.width) * uint32(destinationRect.height)
		rgbCount     = sourceLength * 4
		rgbRaw       = make([]uint8, rgbCount)
	)

	if px.packType == 3 {
		k := 0
		for i := uint32(0); i < sourceLength; i++ {
			rgbRaw[k] = uint8((((pxShortArray[i]) & 0x7C00) >> 10) << 3)
			k++
			rgbRaw[k] = uint8(((pxShortArray[i] & 0x03E0) >> 5) << 3)
			k++
			rgbRaw[k] = uint8((pxShortArray[i] & 0x001F) << 3)
			k++
			rgbRaw[k] = 0xFF // UINT8_MAX
			k++
		}
	} else {
		k := 0
		for i := uint32(0); i < sourceLength; i++ {
			rgbRaw[k] = uint8((pxArray[i] & 0xFF0000) >> 16)
			k++
			rgbRaw[k] = uint8((pxArray[i] & 0xFF00) >> 8)
			k++
			rgbRaw[k] = uint8(pxArray[i] & 0xFF)
			k++
			rgbRaw[k] = uint8((pxArray[i] & 0xFF000000) >> 24)
			k++
		}
	}

	img := image.NewRGBA(image.Rect(int(px.bounds.x), int(px.bounds.y), int(px.bounds.width), int(px.bounds.height)))
	for y := 0; y < int(px.bounds.height); y++ {
		for x := 0; x < int(px.bounds.width); x++ {
			idx := (int(px.bounds.width)*y + x) << 2
			img.Set(x, y, color.RGBA{
				R: rgbRaw[idx+0],
				G: rgbRaw[idx+1],
				B: rgbRaw[idx+2],
				A: rgbRaw[idx+3],
			})
		}
	}

	return img, nil
}
