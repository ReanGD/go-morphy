package std

// StrStr - string + string
type StrStr struct {
	Key   string
	Value string
}

// StrStrs - string + []string
type StrStrs struct {
	Key   string
	Value []string
}

// StrBytes - string + []byte
type StrBytes struct {
	Key   string
	Value []byte
}

// StrUints16 - string + []uint16
type StrUints16 struct {
	Key   string
	Value []uint16
}

// StrUints16Arr string + [][]uint16
type StrUints16Arr struct {
	Key   string
	Value [][]uint16
}

// StrUint32 - string + uint32
type StrUint32 struct {
	Key   string
	Value uint32
}
