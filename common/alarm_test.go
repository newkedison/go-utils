package common_test

import (
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAlarm(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	assert.Equal(a.DataId(), "id")
	assert.Equal(a.Range(), common.NewRange(10, 20))
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 0)
	assert.Equal(a.IsEnabled(), true)
	assert.Equal(a.State(), common.AutoCanceled)
	assert.False(a.IsAlarming())
}

func TestAlarmSetRange(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	assert.Equal(a.Range(), common.NewRange(10, 20))
	a.SetRange(common.NewRange(20, 30))
	assert.Equal(a.Range(), common.NewRange(20, 30))
}

func TestAlarmSetIgnoreCount(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	assert.EqualValues(a.IgnoreCount(), 3)
	a.SetIgnoreCount(0)
	assert.EqualValues(a.IgnoreCount(), 0)
	a.SetIgnoreCount(1e9)
	assert.EqualValues(a.IgnoreCount(), 1e9)
}

func loopCheck(a *common.Alarm, count int) {
	for i := 0; i < count; i++ {
		a.Check()
	}
}

func TestAlarmCheck(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	d.SetValue(5)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 0)
	assert.Equal(a.State(), common.AutoCanceled)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 1)
	assert.Equal(a.State(), common.AutoCanceled)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 2)
	assert.Equal(a.State(), common.AutoCanceled)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 3)
	assert.Equal(a.State(), common.AutoCanceled)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 4)
	assert.Equal(a.State(), common.UnderFlow)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 5)
	assert.Equal(a.State(), common.UnderFlow)
	d.SetValue(15)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 5)
	assert.Equal(a.State(), common.UnderFlow)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 0)
	assert.Equal(a.State(), common.AutoCanceled)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 0)
	assert.Equal(a.State(), common.AutoCanceled)
	d.SetValue(25)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 0)
	assert.Equal(a.State(), common.AutoCanceled)
	loopCheck(&a, 10)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 10)
	assert.Equal(a.State(), common.OverFlow)
	d.SetValue(5)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 10)
	assert.Equal(a.State(), common.OverFlow)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 11)
	assert.Equal(a.State(), common.UnderFlow)
	d.SetValue(25)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 11)
	assert.Equal(a.State(), common.UnderFlow)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 12)
	assert.Equal(a.State(), common.OverFlow)
	a.SetIgnoreCount(100)
	assert.EqualValues(a.IgnoreCount(), 100)
	assert.EqualValues(a.AlarmCount(), 12)
	assert.Equal(a.State(), common.OverFlow)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 100)
	assert.EqualValues(a.AlarmCount(), 13)
	assert.Equal(a.State(), common.OverFlow)
	d.SetValue(15)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 100)
	assert.EqualValues(a.AlarmCount(), 0)
	assert.Equal(a.State(), common.AutoCanceled)
	d.SetValue(5)
	loopCheck(&a, 100)
	assert.EqualValues(a.IgnoreCount(), 100)
	assert.EqualValues(a.AlarmCount(), 100)
	assert.Equal(a.State(), common.AutoCanceled)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 100)
	assert.EqualValues(a.AlarmCount(), 101)
	assert.Equal(a.State(), common.UnderFlow)
	d.SetValue(15)
	a.Check()
	d.SetValue(25)
	loopCheck(&a, 50)
	assert.EqualValues(a.IgnoreCount(), 100)
	assert.EqualValues(a.AlarmCount(), 50)
	assert.Equal(a.State(), common.AutoCanceled)
	a.SetIgnoreCount(30)
	a.Check()
	assert.EqualValues(a.IgnoreCount(), 30)
	assert.EqualValues(a.AlarmCount(), 51)
	assert.Equal(a.State(), common.OverFlow)
}

func TestAlarmCancelAlarm(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	d.SetValue(5)
	loopCheck(&a, 10)
	assert.Equal(a.State(), common.UnderFlow)
	a.CancelAlarm()
	assert.Equal(a.State(), common.ManualCanceled)
	assert.EqualValues(a.IgnoreCount(), 3)
	assert.EqualValues(a.AlarmCount(), 0)
	loopCheck(&a, 10)
	assert.Equal(a.State(), common.UnderFlow)
	assert.EqualValues(a.AlarmCount(), 10)
}

func TestAlarmOnAlarm(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	i := 0
	a.OnAlarm(func(value common.Number, alarm *common.Alarm) {
		i++
		assert.EqualValues(value, 5)
		assert.Equal(&a, alarm)
	})
	d.SetValue(5)
	loopCheck(&a, 4)
	assert.EqualValues(i, 1)
	loopCheck(&a, 100)
	assert.EqualValues(i, 1)
	d.SetValue(15)
	a.Check()
	a.OnAlarm(func(common.Number, *common.Alarm) {
		i *= 10
	})
	a.OnAlarm(func(common.Number, *common.Alarm) {
		i -= 1
	})
	d.SetValue(5)
	loopCheck(&a, 4)
	assert.EqualValues(i, 19)
}

func TestAlarmOnAlarmCanceled(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	i := 0
	a.OnAlarmCanceled(func(value common.Number, alarm *common.Alarm) {
		i++
		assert.True(value == 5 || value == 15)
		assert.Equal(&a, alarm)
	})
	d.SetValue(5)
	loopCheck(&a, 4)
	assert.EqualValues(i, 0)
	loopCheck(&a, 100)
	assert.EqualValues(i, 0)
	d.SetValue(15)
	a.Check()
	assert.EqualValues(i, 1)
	loopCheck(&a, 100)
	assert.EqualValues(i, 1)
	a.OnAlarmCanceled(func(common.Number, *common.Alarm) {
		i *= 10
	})
	a.OnAlarmCanceled(func(common.Number, *common.Alarm) {
		i -= 1
	})
	d.SetValue(5)
	loopCheck(&a, 4)
	d.SetValue(15)
	a.Check()
	assert.EqualValues(i, 19)
	a.OnAlarmCanceled(func(common.Number, *common.Alarm) {
		i -= 1
		assert.EqualValues(a.State(), common.ManualCanceled)
	})
	d.SetValue(5)
	loopCheck(&a, 4)
	a.CancelAlarm()
	assert.EqualValues(i, 198)
	d.SetValue(15)
	a.Check()
	assert.EqualValues(i, 198)
}

func timeNear(t1, t2 time.Time) bool {
	diff := t1.UnixNano() - t2.UnixNano()
	return diff < 100000 && diff > -100000
}

func sleep() {
	time.Sleep(100 * time.Microsecond)
}

func TestAlarmLastValueAndTime(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	d.SetValue(5)
	loopCheck(&a, 10)
	assert.EqualValues(a.LastAlarmValue(), 5)
	assert.True(timeNear(a.LastAlarmTime(), time.Now()))
	a.CancelAlarm()
	assert.EqualValues(a.LastCancelValue(), 5)
	assert.True(timeNear(a.LastCancelTime(), time.Now()))
	loopCheck(&a, 10)
	d.SetValue(15)
	sleep()
	a.Check()
	assert.EqualValues(a.LastCancelValue(), 15)
	assert.True(timeNear(a.LastCancelTime(), time.Now()))
}

func TestAlarmEnable(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	assert.True(a.IsEnabled())
	a.Disable()
	assert.False(a.IsEnabled())
	a.SetEnable(false)
	assert.False(a.IsEnabled())
	d.SetValue(5)
	loopCheck(&a, 1000)
	assert.EqualValues(a.AlarmCount(), 0)
	a.Enable()
	loopCheck(&a, 1000)
	assert.EqualValues(a.AlarmCount(), 1000)
	assert.Equal(a.State(), common.UnderFlow)
	i := 0
	a.OnAlarmCanceled(func(common.Number, *common.Alarm) {
		i++
	})
	assert.False(timeNear(a.LastCancelTime(), time.Now()))
	a.Disable()
	assert.EqualValues(a.AlarmCount(), 0)
	assert.True(timeNear(a.LastCancelTime(), time.Now()))
	assert.Equal(a.State(), common.AutoCanceled)
	assert.EqualValues(i, 1)
}

func TestNewAlarmRecord(t *testing.T) {
	assert := assert.New(t)
	d := common.NewNamedData("id", "name", 99, common.NewRange(0, 100))
	a := common.NewAlarm(&d, common.NewRange(10, 20), 3)
	r := common.NewAlarmRecord(&a)
	assert.Equal(r.DataId, "id")
	assert.Equal(r.State, common.AutoCanceled)
	assert.Equal(r.Value, a.LastCancelValue())
	assert.Equal(r.Time, a.LastCancelTime())
	d.SetValue(5)
	loopCheck(&a, 10)
	r = common.NewAlarmRecord(&a)
	assert.Equal(r.DataId, "id")
	assert.Equal(r.State, common.UnderFlow)
	assert.EqualValues(r.Value, 5)
	assert.True(timeNear(r.Time, time.Now()))
	sleep()
	d.SetValue(25)
	loopCheck(&a, 10)
	r = common.NewAlarmRecord(&a)
	assert.Equal(r.DataId, "id")
	assert.Equal(r.State, common.OverFlow)
	assert.EqualValues(r.Value, 25)
	assert.True(timeNear(r.Time, time.Now()))
}

func TestAlarmRecordMarshalBinary(t *testing.T) {
	assert := assert.New(t)
	r := common.AlarmRecord{
		DataId: "aaa",
		State:  common.UnderFlow,
		Value:  -100,
		Bound:  100,
		Time:   time.Now().Truncate(time.Hour),
	}
	data, err := r.MarshalBinary()
	assert.Nil(err)
	var r2 common.AlarmRecord
	assert.Nil(r2.UnmarshalBinary(data))
	assert.Equal(r2.DataId, "aaa")
	assert.EqualValues(r2.State, common.UnderFlow)
	assert.EqualValues(r2.Value, -100)
	assert.EqualValues(r2.Bound, 100)
	assert.EqualValues(r2.Time, time.Now().Truncate(time.Hour))
}
