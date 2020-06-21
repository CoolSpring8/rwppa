// Copyright (C) 2020  CoolSpring8

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func newLoginWindow() (string, string, string) {
	a := app.New()
	w := a.NewWindow("rwppa")

	loginForm := widget.NewForm()

	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	listenAddrEntry := widget.NewEntry()

	loginForm.Append("Username", usernameEntry)
	loginForm.Append("Password", passwordEntry)
	loginForm.Append("Listen Address", listenAddrEntry)

	loginButton := widget.NewButton("Log in", func() {
		a.Quit()
	})

	w.SetContent(widget.NewVBox(
		loginForm,
		loginButton,
	))
	w.ShowAndRun()

	username := usernameEntry.Text
	password := passwordEntry.Text
	listenAddr := listenAddrEntry.Text

	return username, password, listenAddr
}
