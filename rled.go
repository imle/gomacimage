package gomacimage

import (
	"errors"
	"image"
	"image/color"
)

type RleOpCode uint8

const (
	RleOpCodeEndOfFrame RleOpCode = iota
	RleOpCodeLineStart
	RleOpCodePixelData
	RleOpCodeTransparentRun
	RleOpCodePixelRun
)

func RledFromBytes(b []byte) (sprites []image.Image, err error) {
	parser := dataStructureParse{
		d:   NewBigEndianDataView(b),
		pos: 0,
	}

	// The first part of the rlë resource is the preamble or header. This begins
	// with the dimensions of the sprite.
	width := parser.readWord()
	height := parser.readWord()

	// Following the dimensions is the number of bytes per pixel.
	bitsPerPixel := parser.readWord()

	// There are then two bytes which appear to be unused.
	_ = parser.readDataUint8(2)

	// Followed by the number of frames
	frameCount := parser.readWord()

	// And again there seems to be another run of 6 unused bytes.
	_ = parser.readDataUint8(6)

	// We're going to assume a colour depth of 16. Anything else will trigger an error.
	if bitsPerPixel != 16 {
		return nil, errors.New("invalid color depth in rlëD resource")
	}

	sprites = make([]image.Image, frameCount, frameCount)

	position := uint32(0)
	rowStart := uint32(0)
	currentLine := int32(-1)
	currentColumn := int32(-1)
	offset := int32(0)

	opCode := RleOpCode(0)
	count := uint32(0)
	pixel := uint16(0)
	currentFrame := int32(0)

	var sprite *image.RGBA
	sprite = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for {
		position = uint32(parser.pos)
		if position >= uint32(len(parser.d.buffer)) {
			return nil, errors.New("early end-of-resource encountered in rlëD resource")
		}

		off := (position - rowStart) & 0x03
		if rowStart != 0 && off != 0 {
			position += 4 - off
			parser.pos += int(4 - (count & 0x03))
		}

		count = parser.readDWord()
		opCode = RleOpCode((count & 0xFF000000) >> 24)
		count &= 0x00FFFFFF

		switch opCode {
		case RleOpCodeEndOfFrame:
			if currentLine != int32(height-1) {
				return nil, errors.New("incorrect number of scanlines in rlëD resource")
			}

			sprites[currentFrame] = sprite

			currentFrame++
			if currentFrame >= int32(frameCount) {
				// Finished parsing everything successfully.
				return sprites, nil
			}

			sprite = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
			currentLine = -1

		case RleOpCodeLineStart:
			currentLine++
			currentColumn = 0
			rowStart = uint32(parser.pos)

		case RleOpCodePixelData:
			for i := uint32(0); i < count; i += 2 {
				pixel = parser.readWord()
				offset = currentLine*int32(width) + currentColumn
				writePixelData(sprite, width, height, offset, pixel)
				currentColumn++
			}

			if count&0x03 > 0 {
				parser.pos += int(4 - (count & 0x03))
			}

		case RleOpCodeTransparentRun:
			currentColumn += int32(count >> ((bitsPerPixel >> 3) - 1))

		case RleOpCodePixelRun:
			_ = parser.readDWord()

			for i := uint32(0); i < count; i += 4 {
				offset = currentLine*int32(width) + currentColumn
				writePixelData(sprite, width, height, offset, pixel)
				currentColumn++

				if i+2 < count {
					offset = currentLine*int32(width) + currentColumn
					writePixelData(sprite, width, height, offset, pixel)
					currentColumn++
				}
			}

		default:
			return nil, errors.New("invalid opcode encountered in rlëD resource")
		}
	}
}

func writePixelData(sprite *image.RGBA, width uint16, height uint16, currentOffset int32, col uint16) {
	var blue = col & 0x001F
	var green = (col & 0x03E0) >> 5
	var red = (col & 0x7C00) >> 10
	var alpha = 0xFF

	blue = blue << 3
	green = green << 3
	red = red << 3

	blue |= blue >> 5
	green |= green >> 5
	red |= red >> 5

	// of = y * width + x
	// of - y * width = x
	y := int((currentOffset - currentOffset%int32(sprite.Rect.Size().X)) / int32(sprite.Rect.Size().Y))
	x := int(currentOffset - int32(y)*int32(width))

	sprite.SetRGBA(x, y, color.RGBA{
		R: uint8(0xFF & red),
		G: uint8(0xFF & green),
		B: uint8(0xFF & blue),
		A: uint8(0xFF & alpha),
	})

	//	return rgb + (alpha * 0x01000000);
}
