package condition

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleCodec struct {
	ConditionChoice
}

func NewSubtitleCodecCondition() *SubtitleCodec {
	return &SubtitleCodec{}
}

func (c *SubtitleCodec) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and codec matches
		if s.CodecType == "subtitle" && c.checkString(s.CodecName) {
			return true
		}
	}
	return false
}

func (c *SubtitleCodec) Name() string {
	return "Subtitle codec"
}

func (c *SubtitleCodec) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *SubtitleCodec) New() Condition {
	return &SubtitleCodec{
		ConditionChoice{
			value: c.value,
		},
	}
}

func (c *SubtitleCodec) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (sc *SubtitleCodec) GetEntry() fyne.Widget {
	selectWidget := widget.NewSelect([]string{"subrip", "ass"}, func(s string) {
		sc.value = s
	})
	return selectWidget
}
