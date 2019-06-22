package common_test

import (
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNameData(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 100, common.NewRange(0, 100))
	assert.Equal(d.Id(), "id")
	assert.Equal(d.Name(), "name")
	assert.EqualValues(d.Value(), 100)
	assert.Equal(d.Range(), common.NewRange(0, 100))
}

func TestNameDataSetName(t *testing.T) {
	assert := assert.New(t)
	var d common.NamedData
	d.SetName("aaa")
	assert.Equal(d.Name(), "aaa")
}

func TestNameDataSetRange(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 100, common.NewRange(0, 100))
	assert.EqualValues(d.Value(), 100)
	d.SetRane(common.NewRange(1, 2))
	assert.Equal(d.Range(), common.NewRange(1, 2))
	assert.EqualValues(d.Value(), 2)
	d.SetRane(common.NewRange(10, 100))
	assert.Equal(d.Range(), common.NewRange(10, 100))
	assert.EqualValues(d.Value(), 10)
}

func TestNameDataSetValue(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 100, common.NewRange(0, 100))
	assert.False(d.SetValue(-1))
	assert.False(d.SetValue(101))
	assert.True(d.SetValue(0))
	assert.True(d.SetValue(10))
	assert.True(d.SetValue(100))
}

func TestNameDataAutoCheck(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 100, common.NewRange(0, 100))
	assert.False(d.IsAutoCheck())
	d.SetAutoCheck(true)
	assert.True(d.IsAutoCheck())
	d.SetAutoCheck(false)
	assert.False(d.IsAutoCheck())
	assert.False(d.Warning().IsEnabled())
	assert.False(d.Error().IsEnabled())
	assert.EqualValues(d.Warning().AlarmCount(), 0)
	assert.EqualValues(d.Error().AlarmCount(), 0)
	d.Warning().SetRange(common.NewRange(20, 30))
	d.Warning().SetIgnoreCount(3)
	d.Warning().Enable()
	d.Error().SetRange(common.NewRange(10, 40))
	d.Error().SetIgnoreCount(3)
	d.Error().Enable()
	d.SetValue(5)
	d.SetValue(15)
	d.SetValue(25)
	d.SetValue(35)
	d.SetValue(45)
	assert.EqualValues(d.Warning().AlarmCount(), 0)
	assert.EqualValues(d.Error().AlarmCount(), 0)
}
