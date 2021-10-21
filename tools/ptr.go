package tools

func BoolPtr(in bool) (out *bool) {
	out = new(bool)
	*out = in

	return
}

func StringPtr(in string) (out *string) {
	if in == "" {
		return
	}
	out = new(string)
	*out = in

	return
}

func Int64Ptr(in int64) (out *int64) {
	if in == 0 {
		return
	}
	out = new(int64)
	*out = in

	return
}

func Int32Ptr(in int32) (out *int32) {
	if in == 0 {
		return
	}
	out = new(int32)
	*out = in

	return
}

func Uint64Ptr(in uint64) (out *uint64) {
	if in == 0 {
		return
	}
	out = new(uint64)
	*out = in

	return
}

func Uint32Ptr(in uint32) (out *uint32) {
	if in == 0 {
		return
	}
	out = new(uint32)
	*out = in

	return
}

func IntPtr(in int) (out *int) {
	if in == 0 {
		return
	}
	out = new(int)
	*out = in

	return
}

func UintPtr(in uint) (out *uint) {
	if in == 0 {
		return
	}
	out = new(uint)
	*out = in

	return
}
