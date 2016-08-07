package dawg

// status ok
const (
	// constRoot - Root index
	constRoot          uint32 = 0
	constPrecisionMask uint32 = 0xFFFFFFFF
	constIsLeafBit     uint32 = 1 << 31
	constHasLeafBit    uint32 = 1 << 8
	constExtensionBit  uint32 = 1 << 9
)

// Check if a unit has a leaf as a child or not.
func hasLeaf(base uint32) bool {
	const mask uint32 = constHasLeafBit
	return (base & mask) != 0
}

// Check if a unit corresponds to a leaf or not.
func getValue(base uint32) uint32 {
	const mask uint32 = (^constIsLeafBit) & constPrecisionMask
	return base & mask
}

// Read a label with a leaf flag from a non-leaf unit.
func getLabel(base uint32) uint32 {
	const mask uint32 = constIsLeafBit | 0xFF
	return base & mask
}

// Read an offset to child units from a non-leaf unit.
func getOffset(base uint32) uint32 {
	return ((base >> 10) << ((base & constExtensionBit) >> 6)) & constPrecisionMask
}
