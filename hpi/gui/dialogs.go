package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func ShowInfo(
	parent *gtk.Window,
	title string,
	message string,
) {
	showMessage(
		parent,
		gtk.MESSAGE_INFO,
		title,
		message,
	)
}

func ShowWarning(
	parent *gtk.Window,
	title string,
	message string,
) {
	showMessage(
		parent,
		gtk.MESSAGE_WARNING,
		title,
		message,
	)
}

func ShowError(
	parent *gtk.Window,
	title string,
	message string,
) {
	showMessage(
		parent,
		gtk.MESSAGE_ERROR,
		title,
		message,
	)
}

func AskConfirmation(
	parent *gtk.Window,
	title string,
	message string,
) bool {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_QUESTION,
		gtk.BUTTONS_YES_NO,
		message,
	)

	dialog.SetTitle(title)
	defer dialog.Destroy()

	response := dialog.Run()

	return response == gtk.RESPONSE_YES
}

func showMessage(
	parent *gtk.Window,
	messageType gtk.MessageType,
	title string,
	message string,
) {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		messageType,
		gtk.BUTTONS_OK,
		message,
	)

	dialog.SetTitle(title)
	defer dialog.Destroy()

	dialog.Run()
}
