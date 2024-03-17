package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/core/condition"
)

type ConditionalWidget struct {
	widget.BaseWidget
	key       *widget.Select
	Choice    condition.Condition
	container *fyne.Container
}

var conditions = []condition.Condition{
	condition.NewContainerCondition(),
	condition.NewVideoTitleCondition(),
	condition.NewBitrateCondition(),
	condition.NewAudioLanguageCondition(),
	condition.NewSubtitleLanguageCondition(),
	condition.NewSubtitleCodecCondition(),
}

var choices = make([]string, len(conditions))

func NewConditionalWidget(fv *fyne.Container) *ConditionalWidget {

	for i, c := range conditions {
		choices[i] = c.Name()
	}

	c := &ConditionalWidget{}

	c.key = widget.NewSelect(choices, nil)

	disableSelect := widget.NewSelect([]string{"Select a condition"}, nil)
	disableSelect.Disable()

	disableSelect2 := widget.NewSelect([]string{"Select a value"}, nil)
	disableSelect2.Disable()

	c.container = container.NewGridWithColumns(3, c.key, disableSelect, disableSelect2)

	c.key.OnChanged = func(s string) {
		for _, cond := range conditions {
			if cond.Name() == s {
				c.Choice = cond.New()
				choiceWidget := widget.NewSelect(cond.GetPossibleConditions(), func(s string) {
					c.Choice.SetCondition(s)
				})
				c.container.Objects[1] = choiceWidget
				c.container.Objects[2] = c.Choice.GetEntry()
			}
		}
	}

	c.ExtendBaseWidget(c)

	return c
}

func (c *ConditionalWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.container)
}

func (c *ConditionalWidget) MinSize() fyne.Size {
	return fyne.NewSize(30, 40)
}
