package common

import (
	"github.com/newkedison/go-utils/internal/types"
	"time"
)

type AlarmState int32

const (
	OverFlow       AlarmState = 0
	UnderFlow      AlarmState = 1
	AutoCanceled   AlarmState = 2
	ManualCanceled AlarmState = 3
)

var AlarmState_name = map[int32]string{
	0: "Overflow",
	1: "Underflow",
	2: "AutoCanceled",
	3: "ManualCanceled",
}

var AlarmState_value = map[string]int32{
	"Overflow":       0,
	"Underflow":      1,
	"AutoCanceled":   2,
	"ManualCanceled": 3,
}

type Alarm struct {
	dataId          string
	dataPtr         *Number
	alarmRange      Range
	ignoreCount     uint32
	alarmCount      uint32
	enabled         bool
	state           AlarmState
	lastAlarmValue  Number
	lastCancelValue Number
	lastAlarmTime   time.Time
	lastCancelTime  time.Time
	sigAlarm        SignalAlarm
	sigCanceled     SignalAlarm
}

func NewAlarm(data *NamedData, alarmRange Range, ignoreCount uint32) Alarm {
	return Alarm{
		dataId:      data.id,
		dataPtr:     &data.value,
		alarmRange:  alarmRange,
		ignoreCount: ignoreCount,
		alarmCount:  0,
		enabled:     true,
		state:       AutoCanceled,
	}
}

func (a *Alarm) DataId() string {
	return a.dataId
}

func (a *Alarm) Range() Range {
	return a.alarmRange
}

func (a *Alarm) SetRange(newRange Range) *Alarm {
	a.alarmRange.Change(newRange.low, newRange.high)
	return a
}

func (a *Alarm) OnAlarm(f func(Number, *Alarm)) *Alarm {
	a.sigAlarm.Connect(f)
	return a
}

func (a *Alarm) OnAlarmCanceled(f func(Number, *Alarm)) *Alarm {
	a.sigCanceled.Connect(f)
	return a
}

func (a *Alarm) AlarmCount() uint32 {
	return a.alarmCount
}

func (a *Alarm) IgnoreCount() uint32 {
	return a.ignoreCount
}

func (a *Alarm) SetIgnoreCount(c uint32) *Alarm {
	a.ignoreCount = c
	return a
}

func (a *Alarm) State() AlarmState {
	return a.state
}

func (a *Alarm) IsAlarming() bool {
	return a.state == UnderFlow || a.state == OverFlow
}

func (a *Alarm) Check() {
	if !a.enabled {
		return
	}
	value := *a.dataPtr
	if value < a.alarmRange.low || value > a.alarmRange.high {
		a.alarmCount++
		if a.alarmCount > a.ignoreCount {
			if value < a.alarmRange.low && a.state != UnderFlow {
				a.state = UnderFlow
				a.lastAlarmValue = value
				a.lastAlarmTime = time.Now()
				a.sigAlarm.fire(value, a)
			} else if value > a.alarmRange.high && a.state != OverFlow {
				a.state = OverFlow
				a.lastAlarmValue = value
				a.lastAlarmTime = time.Now()
				a.sigAlarm.fire(value, a)
			}
		}
	} else {
		a.alarmCount = 0
		if a.IsAlarming() {
			a.state = AutoCanceled
			a.lastCancelValue = value
			a.lastCancelTime = time.Now()
			a.sigCanceled.fire(value, a)
		}
	}
}

func (a *Alarm) CancelAlarm() {
	a.alarmCount = 0
	a.state = ManualCanceled
	a.lastCancelValue = *a.dataPtr
	a.lastCancelTime = time.Now()
	a.sigCanceled.fire(*a.dataPtr, a)
}

func (a *Alarm) LastAlarmValue() Number {
	return a.lastAlarmValue
}

func (a *Alarm) LastAlarmTime() time.Time {
	return a.lastAlarmTime
}

func (a *Alarm) LastCancelValue() Number {
	return a.lastCancelValue
}

func (a *Alarm) LastCancelTime() time.Time {
	return a.lastCancelTime
}

func (a *Alarm) IsEnabled() bool {
	return a.enabled
}

func (a *Alarm) SetEnable(en bool) {
	if a.enabled == en {
		return
	}
	a.alarmCount = 0
	if a.enabled && a.IsAlarming() {
		a.state = AutoCanceled
		a.lastCancelValue = *a.dataPtr
		a.lastCancelTime = time.Now()
		a.sigCanceled.fire(a.lastCancelValue, a)
	}
	a.enabled = en
}

func (a *Alarm) Enable() {
	if !a.enabled {
		a.SetEnable(true)
	}
}

func (a *Alarm) Disable() {
	if a.enabled {
		a.SetEnable(false)
	}
}

type AlarmRecord struct {
	DataId string
	State  AlarmState
	Value  Number
	Bound  Number
	Time   time.Time
}

func NewAlarmRecord(a *Alarm) AlarmRecord {
	rcd := new(AlarmRecord)
	rcd.DataId = a.dataId
	rcd.State = a.state
	if a.state == OverFlow || a.state == UnderFlow {
		rcd.Value = a.lastAlarmValue
		rcd.Time = a.lastAlarmTime
		if rcd.State == OverFlow {
			rcd.Bound = a.alarmRange.high
		} else {
			rcd.Bound = a.alarmRange.low
		}
	} else {
		rcd.Value = a.lastCancelValue
		rcd.Bound = 0
		rcd.Time = a.lastCancelTime
	}
	return *rcd
}

func (a *AlarmRecord) ToProtoMessage() *types.SAlarmInfo {
	return &types.SAlarmInfo{
		Id:    a.DataId,
		State: types.AlarmState(uint32(a.State)),
		Value: a.Value.ToProtoMessage(),
		Bound: a.Bound.ToProtoMessage(),
		Time:  a.Time.UnixNano() / 1000000,
	}
}

func (v *AlarmRecord) FromProtoMessage(p *types.SAlarmInfo) {
	v.DataId = p.Id
	v.State = AlarmState(uint32(p.State))
	v.Value.FromProtoMessage(p.Value)
	v.Bound.FromProtoMessage(p.Bound)
	v.Time = time.Unix(0, p.Time*1000000)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a *AlarmRecord) MarshalBinary() (data []byte, err error) {
	defer SetErrorWhenMarshalObjectErrorPanic("common.AlarmRecord", &err)
	return MarshalProtoMessage(a.ToProtoMessage())
}

// UnmarshalBinaryWithSize implements the common.BinaryUnmarshalerWithSize interface.
func (v *AlarmRecord) UnmarshalBinaryWithSize(data []byte) (_ int, err error) {
	defer SetErrorWhenUnmarshalObjectErrorPanic("common.AlarmRecord", &err)()
	var result types.SAlarmInfo
	used := UnmarshalProtoMessage(data, &result)
	v.FromProtoMessage(&result)
	return used, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (v *AlarmRecord) UnmarshalBinary(data []byte) error {
	_, err := v.UnmarshalBinaryWithSize(data)
	return err
}

// vim: fdm=syntax fdn=1
