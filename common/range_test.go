package common_test

import (
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkRange(r common.Range, low common.Number, high common.Number) bool {
	return r.Low() == low && r.High() == high
}

func TestNewRange(t *testing.T) {
	assert := assert.New(t)
	a := common.NewRange(1, 2)
	assert.True(checkRange(a, 1, 2))
	a = common.NewRange(2, 1)
	assert.True(checkRange(a, 1, 2))
	a = common.NewRange(1, 1)
	assert.True(checkRange(a, common.MinNumber, common.MaxNumber))

	assert.EqualValues(a.Low(), common.MinNumber)
	assert.EqualValues(a.High(), common.MaxNumber)
}

func TestRangeChange(t *testing.T) {
	assert := assert.New(t)
	a := common.NewRange(1, 2)
	assert.True(checkRange(a, 1, 2))
	a.Change(1, 2)
	assert.True(checkRange(a, 1, 2))
	a.Change(2, 1)
	assert.True(checkRange(a, 1, 2))
	a.Change(1, 1)
	assert.True(checkRange(a, common.MinNumber, common.MaxNumber))
}
