package common

import (
	"github.com/newkedison/go-utils/internal/types"
)

type Range struct {
	low  Number
	high Number
}

var maxRange Range = Range{MinNumber, MaxNumber}

func MaxRange() Range {
	return maxRange
}

func NewRange(low Number, high Number) Range {
	if low < high {
		return Range{low, high}
	}
	if low == high {
		return Range{MinNumber, MaxNumber}
	}
	return Range{high, low}
}

func convertTwoNumber(a *types.WSNumber, b *types.WSNumber) (Number, Number) {
	var l Number
	var h Number
	l.FromProtoMessage(a)
	h.FromProtoMessage(b)
	return l, h
}

func NewRangeFromInternalType(low *types.WSNumber, high *types.WSNumber) Range {
	return NewRange(convertTwoNumber(low, high))
}

func (r *Range) ChangeFromInternalType(low *types.WSNumber, high *types.WSNumber) {
	r.Change(convertTwoNumber(low, high))
}

func (r *Range) Change(low Number, high Number) {
	if low == high {
		r.low = MinNumber
		r.high = MaxNumber
	} else if low < high {
		r.low = low
		r.high = high
	} else {
		r.low = high
		r.high = low
	}
}

func (r Range) Low() Number {
	return r.low
}

func (r Range) High() Number {
	return r.high
}
