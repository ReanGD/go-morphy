package dawg

// status ok

// Completer ...
type Completer struct {
	dic        *Dictionary
	guide      *Guide
	lastIndex  uint32
	indexStack []uint32
	key        []byte
}

func (c *Completer) value() uint32 {
	return c.dic.value(c.lastIndex)
}

func (c *Completer) start(index uint32, prefix []byte) {
	c.key = prefix

	if c.guide.size() != 0 {
		c.indexStack = []uint32{index}
		c.lastIndex = constRoot
	} else {
		c.indexStack = []uint32{}
	}
}

// Gets the next key
func (c *Completer) next() bool {
	lenIndexStack := len(c.indexStack)
	if lenIndexStack == 0 {
		return false
	}

	var ok bool
	index := c.indexStack[lenIndexStack-1]

	if c.lastIndex != constRoot {
		childLabel := c.guide.child(index) // UCharType

		if childLabel != 0 {
			// Follows a transition to the first child.
			index, ok = c.follow(childLabel, index)
			if !ok {
				return false
			}
		} else {
			for {
				siblingLabel := c.guide.sibling(index)

				// Moves to the previous node.
				if len(c.key) > 0 {
					c.key = c.key[:len(c.key)-1]
				}

				c.indexStack = c.indexStack[:len(c.indexStack)-1]
				if len(c.indexStack) == 0 {
					return false
				}

				index = c.indexStack[len(c.indexStack)-1]
				if siblingLabel != 0 {
					// Follows a transition to the next sibling.
					index, ok = c.follow(siblingLabel, index)
					if !ok {
						return false
					}
					break
				}
			}
		}
	}
	return c.findTerminal(index)
}

func (c *Completer) follow(label uint8, index uint32) (uint32, bool) {
	nextIndex, ok := c.dic.followChar(uint32(label), index)
	if ok {
		c.key = append(c.key, label)
		c.indexStack = append(c.indexStack, nextIndex)
	}

	return nextIndex, ok
}

func (c *Completer) findTerminal(index uint32) bool {
	var ok bool
	for !c.dic.hasValue(index) {
		label := c.guide.child(index)

		index, ok = c.dic.followChar(uint32(label), index)
		if !ok {
			return false
		}

		c.key = append(c.key, label)
		c.indexStack = append(c.indexStack, index)
	}

	c.lastIndex = index
	return true
}

// NewCompleter - constructor for Completer
func NewCompleter(dic *Dictionary, guide *Guide) *Completer {
	return &Completer{dic: dic, guide: guide, lastIndex: 0}
}
