//gui/progress.go
package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type Progress struct {
	dialog *gtk.Dialog
	bar    *gtk.ProgressBar
	label  *gtk.Label
}

func NewProgress(parent *gtk.Window, title string) (*Progress, error) {
	dialog, err := gtk.DialogNew()
	if err != nil {
		return nil, fmt.Errorf("create progress dialog: %w", err)
	}

	dialog.SetTitle(title)
	dialog.SetTransientFor(parent)
	dialog.SetModal(true)
	dialog.SetDefaultSize(460, 160)
	dialog.SetPosition(gtk.WIN_POS_CENTER)

	content, err := dialog.GetContentArea()
	if err != nil {
		dialog.Destroy()
		return nil, fmt.Errorf("get dialog content area: %w", err)
	}

	content.SetBorderWidth(20)

	label, err := gtk.LabelNew("Preparing...")
	if err != nil {
		dialog.Destroy()
		return nil, fmt.Errorf("create progress label: %w", err)
	}

	label.SetHAlign(gtk.ALIGN_START)

	content.PackStart(label, false, false, 0)

	bar, err := gtk.ProgressBarNew()
	if err != nil {
		dialog.Destroy()
		return nil, fmt.Errorf("create progress bar: %w", err)
	}

	bar.SetShowText(true)

	content.PackStart(bar, false, false, 15)

	return &Progress{
		dialog: dialog,
		bar:    bar,
		label:  label,
	}, nil
}

func (p *Progress) SetMessage(message string) {
	if message == "" {
		return
	}

	p.label.SetText(message)
}

func (p *Progress) SetFraction(fraction float64) {
	if fraction < 0 {
		fraction = 0
	}

	if fraction > 1 {
		fraction = 1
	}

	p.bar.SetFraction(fraction)
}

func (p *Progress) SetIndeterminate(indeterminate bool) {
	p.bar.SetPulseStep(0.1)

	if indeterminate {
		p.bar.Pulse()
	}
}

func (p *Progress) Show() {
	p.dialog.ShowAll()
}

func (p *Progress) Close() {
	p.dialog.Destroy()
}
