package gomacimage

const (
	WordSize = 2
)

type macRectangle struct {
	y1 uint16
	x1 uint16
	y2 uint16
	x2 uint16
}

type dataStructureParse struct {
	d      *DataView
	pos    int
	xRatio uint16
	yRatio uint16
}

type regionRect struct {
	x      uint16
	y      uint16
	width  uint16
	height uint16
}

type pixMap struct {
	baseAddress uint32
	rowBytes    uint16
	bounds      regionRect
	pmVersion   uint16
	packType    uint16
	packSize    uint32
	hRes        uint32
	vRes        uint32
	pixelType   uint16
	pixelSize   uint16
	cmpCount    uint16
	cmpSize     uint16
	planeBytes  uint32
	pmTable     uint32
	pmReserved  uint32
}

type bitMap struct {
	baseAddress uint32
	rowBytes    uint16
	bounds      regionRect
}

type colorTable struct {
	seed  uint32
	flags uint16
	size  uint16
	data  []colorRow
}

func (p *dataStructureParse) readDataUint8(len int) []byte {
	d := p.d.buffer[p.pos : p.pos+len : p.pos+len]
	p.pos += len
	return d
}

func (p *dataStructureParse) readData(len int) *DataView {
	var data = NewBigEndianDataView(p.d.buffer[p.pos : p.pos+len : p.pos+len])
	p.pos += len
	return data
}

func (p *dataStructureParse) parsePixMap() pixMap {
	return pixMap{
		baseAddress: p.readDWord(),
		rowBytes:    p.readWord() & 0x7FFF,

		bounds: p.readWHRect(),

		pmVersion: p.readWord(),
		packType:  p.readWord(),
		packSize:  p.readDWord(),

		hRes: p.readFixedPoint(),
		vRes: p.readFixedPoint(),

		pixelType: p.readWord(),
		pixelSize: p.readWord(),
		cmpCount:  p.readWord(),
		cmpSize:   p.readWord(),

		planeBytes: p.readDWord(),
		pmTable:    p.readDWord(),
		pmReserved: p.readDWord(),
	}
}

func (p *dataStructureParse) parseBitMap() bitMap {
	return bitMap{
		baseAddress: p.readDWord(),
		rowBytes:    p.readWord() & 0x7FFF,
		bounds:      p.readWHRect(),
	}
}

type colorRow struct {
	r, g, b, value uint16
}

func (p *dataStructureParse) parseColorTable() colorTable {
	ct := colorTable{
		seed:  p.readDWord(),
		flags: p.readWord(),
		size:  p.readWord() + 1,
	}

	ct.data = make([]colorRow, ct.size)

	for i := uint16(0); i < ct.size; i++ {
		ct.data[i] = colorRow{
			value: p.readWord(),
			r:     p.readWord(),
			g:     p.readWord(),
			b:     p.readWord(),
		}
	}

	return ct
}

func (p *dataStructureParse) readQDRect() macRectangle {
	var rect = macRectangle{
		y1: p.d.GetUint16(p.pos + 0*WordSize),
		x1: p.d.GetUint16(p.pos + 1*WordSize),
		y2: p.d.GetUint16(p.pos + 2*WordSize),
		x2: p.d.GetUint16(p.pos + 3*WordSize),
	}
	p.pos += WordSize * 4
	return rect
}

func (p *dataStructureParse) readWHRect() regionRect {
	var r = p.readQDRect()
	return regionRect{
		x:      r.x1,
		y:      r.y1,
		width:  r.x2 - r.x1,
		height: r.y2 - r.y1,
	}
}

func (p *dataStructureParse) readFixedPoint() uint32 {
	var point = p.d.GetUint32(p.pos) / (1 << 16)
	p.pos += 4
	return point
}

func (p *dataStructureParse) readByte() uint8 {
	var b = p.d.GetUint8(p.pos)
	p.pos++
	return b
}

func (p *dataStructureParse) readDWord() uint32 {
	var word = p.d.GetUint32(p.pos)
	p.pos += 4
	return word
}

func (p *dataStructureParse) readWord() uint16 {
	var word = p.d.GetUint16(p.pos)
	p.pos += 2
	return word
}

func (p *dataStructureParse) readOpCode() uint16 {
	p.pos += p.pos % 2
	return p.readWord()
}
