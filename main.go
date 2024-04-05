package main

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Login Form")
	myWindow.Resize(fyne.NewSize(800, 600))

	UserEntry := widget.NewEntry()
	PassEntry := widget.NewPasswordEntry()
	//textArea := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Username", Widget: UserEntry},
			{Text: "Password", Widget: PassEntry}},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Username:", UserEntry.Text)
			log.Println("Password:", PassEntry.Text)
			myWindow.Close()
		},
	}

	myWindow.Padded()
	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}
