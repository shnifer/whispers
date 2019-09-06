package rid

import "math"

type U16 uint16
type U32 uint32
type U64 uint64

func (i U16) Less(j U16) bool {
	if i > math.MaxUint16/4*3 && j < math.MaxUint16/4 {
		return j < i
	} else {
		return i < j
	}
}

func (i U32) Less(j U32) bool {
	if i > math.MaxUint32/4*3 && j < math.MaxUint32/4 {
		return j < i
	} else {
		return i < j
	}
}

func (i U64) Less(j U64) bool {
	if i > math.MaxUint64/4*3 && j < math.MaxUint64/4 {
		return j < i
	} else {
		return i < j
	}
}
