package gomacimage

import (
	"encoding/binary"
	"math"
)

type DataView struct {
	endian binary.ByteOrder
	buffer []byte
}

func (d DataView) GetLength() int {
	return len(d.buffer)
}

func (d DataView) GetFloat32(byteOffset int) float32 {
	return math.Float32frombits(d.endian.Uint32(d.buffer[byteOffset:]))
}

func (d DataView) GetFloat64(byteOffset int) float64 {
	return math.Float64frombits(d.endian.Uint64(d.buffer[byteOffset:]))
}

func (d DataView) GetInt8(byteOffset int) int8 {
	return int8(d.buffer[byteOffset])
}

func (d DataView) GetInt16(byteOffset int) int16 {
	return int16(d.endian.Uint16(d.buffer[byteOffset:]))
}

func (d DataView) GetInt32(byteOffset int) int32 {
	return int32(d.endian.Uint32(d.buffer[byteOffset:]))
}

func (d DataView) GetUint8(byteOffset int) uint8 {
	return d.buffer[byteOffset]
}

func (d DataView) GetUint16(byteOffset int) uint16 {
	return d.endian.Uint16(d.buffer[byteOffset:])
}

func (d DataView) GetUint32(byteOffset int) uint32 {
	return d.endian.Uint32(d.buffer[byteOffset:])
}

func (d DataView) SetFloat32(byteOffset int, value float32) {
	panic("implement me")
}

func (d DataView) SetFloat64(byteOffset int, value float64) {
	panic("implement me")
}

func (d DataView) SetInt8(byteOffset int, value int8) {
	panic("implement me")
}

func (d DataView) SetInt16(byteOffset int, value int16) {
	panic("implement me")
}

func (d DataView) SetInt32(byteOffset int, value int32) {
	panic("implement me")
}

func (d DataView) SetUint8(byteOffset int, value uint8) {
	panic("implement me")
}

func (d DataView) SetUint16(byteOffset int, value uint16) {
	panic("implement me")
}

func (d DataView) SetUint32(byteOffset int, value uint32) {
	panic("implement me")
}

func NewBigEndianDataView(buffer []byte) *DataView {
	return &DataView{buffer: buffer, endian: binary.BigEndian}
}

func NewLittleEndianDataView(buffer []byte) *DataView {
	return &DataView{buffer: buffer, endian: binary.LittleEndian}
}
