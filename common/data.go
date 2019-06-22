package common

import (
	"errors"
	"github.com/newkedison/go-utils/internal/types"
)

const (
	InvalidValue Number = -1
)

type NamedData struct {
	id               string
	name             string
	value            Number
	dataRange        Range
	warning          Alarm
	fault            Alarm
	sigModified      SignalDataModified
	sigRangeModified SignalDataRangeModified
	sigCheckRead     SignalDataCheck
	sigCheckWrite    SignalDataCheck
	autoCheck        bool
}

func NewNamedData(id string, name string, initValue Number,
	dataRange Range) NamedData {
	d := NamedData{
		id:        id,
		name:      name,
		value:     initValue,
		dataRange: dataRange,
	}
	d.warning = NewAlarm(&d, NewRange(0, 0), 0)
	d.warning.Disable()
	d.fault = NewAlarm(&d, NewRange(0, 0), 0)
	d.fault.Disable()
	return d
}

func (data *NamedData) Id() string {
	return data.id
}

func (data *NamedData) Name() string {
	return data.name
}

func (data *NamedData) SetName(s string) {
	data.name = s
}

func (data *NamedData) Range() Range {
	return data.dataRange
}

func (data *NamedData) SetRane(r Range) bool {
	if !data.IsWritable(data.value) {
		return false
	}
	origin := data.dataRange
	data.dataRange.Change(r.low, r.high)
	data.sigRangeModified.fire(origin, r)
	oldValue := data.value
	if data.value < r.low {
		data.value = r.low
	} else if data.value > r.high {
		data.value = r.high
	}
	if oldValue != data.value {
		data.sigModified.fire(data, oldValue, data.value)
	}
	return true
}

func (data *NamedData) Value() Number {
	if data.IsReadable() {
		return data.value
	} else {
		return InvalidValue
	}
}

func (data *NamedData) SetValue(newValue Number) bool {
	if !data.IsWritable(newValue) {
		return false
	}
	if newValue < data.dataRange.low || newValue > data.dataRange.high {
		return false
	}
	tmp := data.value
	data.value = newValue
	data.sigModified.fire(data, tmp, newValue)
	if data.autoCheck {
		data.warning.Check()
		data.fault.Check()
	}
	return true
}

func (data *NamedData) IsAutoCheck() bool {
	return data.autoCheck
}

func (data *NamedData) SetAutoCheck(b bool) {
	data.autoCheck = b
}

func (data *NamedData) IsReadable() bool {
	return data.sigCheckRead.fire(data, data.value)
}

func (data *NamedData) IsWritable(newValue Number) bool {
	return data.sigCheckWrite.fire(data, newValue)
}

func (data *NamedData) OnModified(
	f func(*NamedData, Number, Number)) {
	data.sigModified.Connect(f)
}

func (data *NamedData) OnRangeModified(f func(Range, Range)) {
	data.sigRangeModified.Connect(f)
}

func (data *NamedData) AddCheckReadMethod(
	f func(*NamedData, Number) bool) {
	data.sigCheckRead.Connect(f)
}

func (data *NamedData) AddCheckWriteMethod(
	f func(*NamedData, Number) bool) {
	data.sigCheckWrite.Connect(f)
}

func (data *NamedData) SignalCheckRead() SignalDataCheck {
	return data.sigCheckRead
}

func (data *NamedData) SignalCheckWrite() SignalDataCheck {
	return data.sigCheckWrite
}

func (data *NamedData) Warning() *Alarm {
	return &data.warning
}

func (data *NamedData) Error() *Alarm {
	return &data.fault
}

func (v *NamedData) ToProtoMessage() *types.WSData {
	return &types.WSData{
		Id:                 v.id,
		Name:               v.name,
		Value:              v.value.ToProtoMessage(),
		RangeLow:           v.dataRange.low.ToProtoMessage(),
		RangeHigh:          v.dataRange.high.ToProtoMessage(),
		WarningLow:         v.warning.alarmRange.low.ToProtoMessage(),
		WarningHigh:        v.warning.alarmRange.high.ToProtoMessage(),
		WarningIgnoreCount: v.warning.ignoreCount,
		ErrorLow:           v.fault.alarmRange.low.ToProtoMessage(),
		ErrorHigh:          v.fault.alarmRange.high.ToProtoMessage(),
		ErrorIgnoreCount:   v.fault.ignoreCount,
	}
}

func (v *NamedData) FromProtoMessage(p *types.WSData) {
	v.id = p.Id
	v.name = p.Name
	v.value.FromProtoMessage(p.Value)
	v.dataRange.ChangeFromInternalType(p.RangeLow, p.RangeHigh)
	v.warning = NewAlarm(v, NewRangeFromInternalType(p.WarningLow, p.WarningHigh), p.WarningIgnoreCount)
	v.fault = NewAlarm(v, NewRangeFromInternalType(p.ErrorLow, p.ErrorHigh), p.ErrorIgnoreCount)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (data *NamedData) MarshalBinary() (result []byte, err error) {
	defer SetErrorWhenMarshalObjectErrorPanic("common.NamedData", &err)()
	return MarshalProtoMessage(data.ToProtoMessage())
}

// UnmarshalBinaryWithSize implements the common.BinaryUnmarshalerWithSize interface.
func (v *NamedData) UnmarshalBinaryWithSize(data []byte) (_ int, err error) {
	defer SetErrorWhenUnmarshalObjectErrorPanic("common.NamedData", &err)()
	if !v.IsWritable(v.value) {
		panic(NewUnmarshalObjectError(errors.New("Unmarshal " + v.id + " fail: not writable.")))
	}
	var result types.WSData
	used := UnmarshalProtoMessage(data, &result)
	v.FromProtoMessage(&result)
	return used, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (v *NamedData) UnmarshalBinary(data []byte) error {
	_, err := v.UnmarshalBinaryWithSize(data)
	return err
}

// vim: fdm=syntax fdn=1
