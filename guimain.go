package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"gitlab.com/tsuchinaga/go-fyne-learning/theme"
)

var (
	w           fyne.Window
	currentPage int
)

func guimain() {
	a := app.New()
	a.Settings().SetTheme(&theme.KoruriTheme{})
	left := widget.NewVBox()
	for i := 0; i < 100; i++ {
		left.Append(widget.NewLabel("foo"))
	}
	w = a.NewWindow("Hoppii info")
	// entry

	// radio, check and select

	box3 := widget.NewVBox()

	w.SetContent(widget.NewVBox(box3))
	assignmentPage()
	w.Resize(fyne.NewSize(300, 512))
	w.ShowAndRun()
}

func NewPage(i int) fyne.CanvasObject {

	moveButtonsBox := widget.NewHBox(
		widget.NewButton("　おしらせ　", announcePage),

		widget.NewButton("　　課題　　", assignmentPage),

		widget.NewButton("ID登録", registerPage),
	)
	moveButtonsBox.CreateRenderer()
	page := widget.NewVBox()
	page.Append(moveButtonsBox)
	registerID(page)
	page = assignmentPageWidget(page)
	page = announceInfoWidget(page)

	return page
}

func assignmentPage() {
	currentPage = 0
	w.SetContent(NewPage(currentPage))
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for range ticker.C {
			// 課題情報を更新する処理
			w.SetContent(NewPage(currentPage))
		}
	}()
}

func announcePage() {

	currentPage = 1

	w.SetContent(NewPage(currentPage))
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for range ticker.C {
			// 課題情報を更新する処理
			w.SetContent(NewPage(currentPage))
		}
	}()
}

func registerPage() {
	currentPage = 2

	w.SetContent(NewPage(currentPage))
}

func registerID(page *widget.Box) {
	if currentPage == 2 {
		inputID := widget.NewEntry()
		inputpass := widget.NewPasswordEntry()

		box2 := widget.NewVBox(
			widget.NewForm(
				widget.NewFormItem("ID", inputID),
				widget.NewFormItem("PASS", inputpass),
			),
			widget.NewButton("登録", func() {

				fmt.Println(inputID.Text, inputpass.Text)
				err := os.Remove("userInfo.csv")
				if err != nil {
					fmt.Println(err)
					return
				}
				file, err := os.OpenFile("userInfo.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer file.Close()

				writer := csv.NewWriter(file)
				defer writer.Flush()

				record := []string{inputID.Text, inputpass.Text}
				err = writer.Write(record)
				if err != nil {
					fmt.Println(err)
					return
				}
			}))

		page.Append(box2)
	}
}

func assignmentPageWidget(page *widget.Box) *widget.Box {
	if currentPage == 0 {
		classInfos := assigntmentInfo()
		assignmentBox := widget.NewVBox()
		if classInfos == nil {

			assignmentBox.Append(widget.NewLabel("no Infomation or network error"))
		} else {
			for _, v := range classInfos {
				fmt.Println(v)
			}
		}
		// 課題を表示するウィジェットを作成

		for _, v := range classInfos {
			assignmentBox.Append(widget.NewLabel(fmt.Sprintf("%s   %d年%d月%d日%d時%d分%d秒", v.name, v.time.Year, v.time.Month, v.time.Day, v.time.Hour, v.time.Min, v.time.Sec)))
		}
		// スクロール可能なウィジェットを作成
		scrollableAssignmentBox := widget.NewScrollContainer(assignmentBox)
		scrollableAssignmentBox.SetMinSize(fyne.NewSize(0, 512)) // 最小サイズを指定する
		page.Append(scrollableAssignmentBox)
	}
	return page
}
func announceInfoWidget(page *widget.Box) *widget.Box {
	if currentPage == 1 {
		announceInfos := announcementInfo()
		announceBox := widget.NewVBox()
		if announceInfos == nil {
			announceBox.Append(widget.NewLabel("no Infomation or network error"))
		} else {
			for _, v := range announceInfos {
				fmt.Println(v)
			}
		}
		// 課題を表示するウィジェットを作成

		for _, v := range announceInfos {
			announceBox.Append(widget.NewLabel(fmt.Sprintf("%s : %s", v.name, v.content)))
		}
		// スクロール可能なウィジェットを作成
		scrollableAssignmentBox := widget.NewScrollContainer(announceBox)
		scrollableAssignmentBox.SetMinSize(fyne.NewSize(0, 512)) // 最小サイズを指定する
		page.Append(scrollableAssignmentBox)
	}
	return page
}
