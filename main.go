package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/helwan-linux/hpi/gui"
	"github.com/helwan-linux/hpi/i18n"
)

func main() {
	gtk.Init(&os.Args)

	window, err := gui.NewWindow()
	if err != nil {
		log.Fatalf("failed to create HPI window: %v", err)
	}

	if len(os.Args) > 1 {
		packagePath := os.Args[1]

		if err := window.OpenPackage(packagePath); err != nil {
			tr := i18n.New()

			window.ShowError(
				tr.Get("package_error"),
				err.Error(),
			)
		}
	}

	window.ShowAll()

	gtk.Main()
}
