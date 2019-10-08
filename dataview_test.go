package gomacimage

import (
	"testing"
)

func TestDataView_GetFloat32(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x10, 0x10, 0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          float32
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          2.8411367e-29,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          2.841137e-29,
			errorExpected: false,
		},
		{
			name:          "out of range",
			byteOffset:    2,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetFloat32() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetFloat32(tt.byteOffset); got != tt.want {
				t.Errorf("GetFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetFloat64(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          float64
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          2.586563270614692e-231,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          2.5865632706146925e-231,
			errorExpected: false,
		},
		{
			name:          "out of range",
			byteOffset:    2,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetFloat64() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetFloat64(tt.byteOffset); got != tt.want {
				t.Errorf("GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetInt8(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          int8
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          0x10,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          0x11,
			errorExpected: false,
		},
		{
			name:          "out of range offset",
			byteOffset:    2,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetInt8() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetInt8(tt.byteOffset); got != tt.want {
				t.Errorf("GetInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetInt16(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          int16
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          0x1010,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          0x1011,
			errorExpected: false,
		},
		{
			name:          "out of range offset",
			byteOffset:    4,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetInt16() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetInt16(tt.byteOffset); got != tt.want {
				t.Errorf("GetInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetInt32(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x10, 0x10, 0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          int32
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          0x10101010,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          0x10101011,
			errorExpected: false,
		},
		{
			name:          "out of range offset",
			byteOffset:    4,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetInt32() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetInt32(tt.byteOffset); got != tt.want {
				t.Errorf("GetInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetUint8(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          uint8
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          0x10,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          0x11,
			errorExpected: false,
		},
		{
			name:          "out of range offset",
			byteOffset:    2,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetUint8() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetUint8(tt.byteOffset); got != tt.want {
				t.Errorf("GetUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetUint16(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          uint16
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          0x1010,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          0x1011,
			errorExpected: false,
		},
		{
			name:          "out of range offset",
			byteOffset:    4,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetUint16() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetUint16(tt.byteOffset); got != tt.want {
				t.Errorf("GetUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetUint32(t *testing.T) {
	d := NewBigEndianDataView([]byte{0x10, 0x10, 0x10, 0x10, 0x11})

	tests := []struct {
		name          string
		byteOffset    int
		want          uint32
		errorExpected bool
	}{
		{
			name:          "0 offset",
			byteOffset:    0,
			want:          0x10101010,
			errorExpected: false,
		},
		{
			name:          "1 offset",
			byteOffset:    1,
			want:          0x10101011,
			errorExpected: false,
		},
		{
			name:          "out of range offset",
			byteOffset:    4,
			want:          0,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.errorExpected {
						t.Error("GetUint32() errored unexpectedly: ", r)
					}
				}
			}()

			if got := d.GetUint32(tt.byteOffset); got != tt.want {
				t.Errorf("GetUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataView_GetLength(t *testing.T) {
	tests := []struct {
		name string
		dv   *DataView
		want int
	}{
		{
			name: "2",
			dv:   NewBigEndianDataView([]byte{0x10, 0x11}),
			want: 2,
		},
		{
			name: "0",
			dv:   NewBigEndianDataView([]byte{}),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dv.GetLength(); got != tt.want {
				t.Errorf("GetLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
