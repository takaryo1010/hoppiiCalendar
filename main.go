package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"gitlab.com/tsuchinaga/go-fyne-learning/theme"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&theme.KoruriTheme{})

	w = a.NewWindow("Hoppii info")
	assignmentPage()
	w.Resize(fyne.NewSize(300, 512))
	w.ShowAndRun()
}
