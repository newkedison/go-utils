package common

type SignalAlarm []func(Number, *Alarm)
type SignalDataCheck []func(*NamedData, Number) bool
type SignalDataModified []func(*NamedData, Number, Number)
type SignalDataRangeModified []func(Range, Range)

func (sig *SignalAlarm) Connect(f func(Number, *Alarm)) {
	(*sig) = append(*sig, f)
}

func (sig *SignalDataCheck) Connect(f func(*NamedData, Number) bool) {
	(*sig) = append(*sig, f)
}

func (sig *SignalDataModified) Connect(
	f func(*NamedData, Number, Number)) {
	(*sig) = append(*sig, f)
}

func (sig *SignalDataRangeModified) Connect(f func(Range, Range)) {
	(*sig) = append(*sig, f)
}

func (sig *SignalAlarm) fire(value Number, alarm *Alarm) {
	for _, f := range *sig {
		f(value, alarm)
	}
}

func (sig *SignalDataCheck) fire(data *NamedData, setValue Number) bool {
	for _, f := range *sig {
		if !f(data, setValue) {
			return false
		}
	}
	return true
}

func (sig *SignalDataModified) fire(data *NamedData, oldValue Number,
	newValue Number) {
	for _, f := range *sig {
		f(data, oldValue, newValue)
	}
}

func (sig *SignalDataRangeModified) fire(oldRange Range, newRange Range) {
	for _, f := range *sig {
		f(oldRange, newRange)
	}
}
