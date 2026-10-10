/*
 * Copyright (C) 2026 Saeed Badreldin
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
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
