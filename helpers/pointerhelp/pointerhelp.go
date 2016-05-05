package pointerhelp

func StringPtr(in string) *string {
	return &in
}

func BoolPtr(in bool) *bool {
	return &in
}

func IntPtr(in int) *int {
	return &in
}

func Int8Ptr(in int8) *int8 {
	return &in
}

func Int16Ptr(in int16) *int16 {
	return &in
}

func Int32Ptr(in int32) *int32 {
	return &in
}

func Int64Ptr(in int64) *int64 {
	return &in
}

func UintPtr(in uint) *uint {
	return &in
}

func Uint8Ptr(in uint8) *uint8 {
	return &in
}

func Uint16Ptr(in uint16) *uint16 {
	return &in
}

func Uint32Ptr(in uint32) *uint32 {
	return &in
}

func Uint64Ptr(in uint64) *uint64 {
	return &in
}

func Float32Ptr(in float32) *float32 {
	return &in
}

func Float64Ptr(in float64) *float64 {
	return &in
}

func Complex64Ptr(in float64) *float64 {
	return &in
}

func Complex128Ptr(in complex128) *complex128 {
	return &in
}

func SafeString(in *string) string {
	if in == nil {
		in = new(string)
	}
	return *in
}

func SafeBool(in *bool) bool {
	if in == nil {
		in = new(bool)
	}
	return *in
}

func SafeInt(in *int) int {
	if in == nil {
		in = new(int)
	}
	return *in
}

func SafeInt8(in *int8) int8 {
	if in == nil {
		in = new(int8)
	}
	return *in
}

func SafeInt16(in *int16) int16 {
	if in == nil {
		in = new(int16)
	}
	return *in
}

func SafeInt32(in *int32) int32 {
	if in == nil {
		in = new(int32)
	}
	return *in
}

func SafeInt64(in *int64) int64 {
	if in == nil {
		in = new(int64)
	}
	return *in
}

func SafeUint(in *uint) uint {
	if in == nil {
		in = new(uint)
	}
	return *in
}

func SafeUint8(in *uint8) uint8 {
	if in == nil {
		in = new(uint8)
	}
	return *in
}

func SafeUint16(in *uint16) uint16 {
	if in == nil {
		in = new(uint16)
	}
	return *in
}

func SafeUint32(in *uint32) uint32 {
	if in == nil {
		in = new(uint32)
	}
	return *in
}

func SafeUint64(in *uint64) uint64 {
	if in == nil {
		in = new(uint64)
	}
	return *in
}

func SafeFloat32(in *float32) float32 {
	if in == nil {
	in = new(float32)
	}
	return *in
}

func SafeFloat64(in *float64) float64 {
	if in == nil {
		in = new(float64)
	}
	return *in
}

func SafeComplex64(in *complex64) complex64 {
	if in == nil {
		in = new(complex64)
	}
	return *in
}

func SafeComplex128(in *complex128) complex128 {
	if in == nil {
		in = new(complex128)
	}
	return *in
}
