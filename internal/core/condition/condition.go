package condition

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type Condition interface {
	Name() string
	Check(data *ffprobe.ProbeData) bool
	GetPossibleConditions() []string
	New() Condition
	SetCondition(condition string)
	GetEntry() fyne.Widget
}

const (
	ContainerConditionName        = "Container"
	VideoTitleConditionName       = "Video Title"
	BitrateConditionName          = "Bitrate"
	AudioLanguageConditionName    = "Audio Language"
	SubtitleLanguageConditionName = "Subtitle Language"
	SubtitleCodecConditionName    = "Subtitle Codec"
)

type ConditionChoice struct {
	condition *ConditionString
	value     string
}

type ConditionString string

var (
	equals          ConditionString = "equals"
	contains        ConditionString = "contains"
	notEquals       ConditionString = "not equals"
	greaterThan     ConditionString = "greater than"
	lessThan        ConditionString = "less than"
	greaterOrEquals ConditionString = "greater or equals"
	lessOrEquals    ConditionString = "less or equals"
)

func FromString(s string) *ConditionString {
	condition := ConditionString(s)
	return &condition
}

func (c *ConditionChoice) checkString(value string) bool {
	switch *c.condition {
	case equals:
		return value == c.value
	case contains:
		return strings.Contains(value, c.value)
	case notEquals:
		return value != c.value
	default:
		return false
	}
}

func (c *ConditionChoice) checkInt(value string) bool {
	intConditionValue, err := strconv.Atoi(c.value)
	if err != nil {
		return false
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return false
	}

	switch *c.condition {
	case equals:
		return intValue == intConditionValue
	case greaterThan:
		return intValue > intConditionValue
	case lessThan:
		return intValue < intConditionValue
	case greaterOrEquals:
		return intValue >= intConditionValue
	case lessOrEquals:
		return intValue <= intConditionValue
	default:
		return false
	}
}

func (c *ConditionChoice) SetValue(value string) {
	c.value = value
}

func (c *ConditionChoice) GetValue() string {
	return c.value
}
