package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

var (
	w           fyne.Window
	currentPage int
)

func NewPage(i int) fyne.CanvasObject {
	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyF5 {
			w.SetContent(loadPage())
			ticker := time.NewTicker(time.Minute * 5)
			go func() {
				if a := NewPage(currentPage); a != nil {
					w.SetContent(a)
				}

				for range ticker.C {

					// 課題情報を更新する処理
					w.SetContent(NewPage(currentPage))
				}
			}()
			// call your function here
		}
	})
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
	if i == currentPage {
		return page
	}
	return nil
}
func loadPage() fyne.CanvasObject {

	moveButtonsBox := widget.NewHBox(
		widget.NewButton("　おしらせ　", announcePage),

		widget.NewButton("　　課題　　", assignmentPage),

		widget.NewButton("ID登録", registerPage),
	)
	moveButtonsBox.CreateRenderer()
	page := widget.NewVBox()
	page.Append(moveButtonsBox)
	loadBox := widget.NewHBox(widget.NewLabel("NowLoading..."))
	page.Append(widget.NewProgressBarInfinite())
	page.Append(loadBox)

	return page
}

func assignmentPage() {

	currentPage = 0

	w.SetContent(loadPage())
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		if a := NewPage(currentPage); a != nil {
			w.SetContent(a)
		}

		for range ticker.C {

			// 課題情報を更新する処理
			w.SetContent(NewPage(currentPage))
		}
	}()
}

func announcePage() {

	currentPage = 1

	w.SetContent(loadPage())
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		if a := NewPage(currentPage); a != nil {
			w.SetContent(a)
		}

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
func hide_particular_challenge() fyne.CanvasObject {
	moveButtonsBox := widget.NewHBox(
		widget.NewButton("　おしらせ　", announcePage),
		widget.NewButton("　　課題　　", assignmentPage),
		widget.NewButton("ID登録", registerPage),
	)
	moveButtonsBox.CreateRenderer()

	// テキスト入力フィールドと追加ボタンの作成
	entry := widget.NewEntry()
	entry.MinSize()
	addButton := widget.NewButton("追加", func() {
		// 入力されたテキストを取得して、remove.csvに追加する
		file, _ := os.OpenFile("remove.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		defer file.Close()
		_, err := file.WriteString(entry.Text + "\n")
		if err != nil {
			fmt.Println("Failed to write to file:", err)
		}
		w.SetContent(hide_particular_challenge())
	})

	// remove.csvから要素を読み込んで、ラベルとしてページに追加する
	file, err := os.Open("remove.csv")
	if err != nil {
		fmt.Println("Failed to open file:", err)
	}
	defer file.Close()

	// ページの作成
	page := widget.NewVBox()
	page.Append(moveButtonsBox)
	page.Append(widget.NewLabel("非表示にしたいキーワードを入力してください"))
	page.Append(widget.NewHBox(entry, addButton))
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		// ラベルと削除ボタンを作成する
		label := widget.NewLabel(text)
		removeButton := widget.NewButton("削除", func() {
			// remove.csvから対応した文字列を削除する
			file, _ := os.OpenFile("remove.csv", os.O_RDWR, 0644)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var lines []string
			for scanner.Scan() {
				if scanner.Text() != text {
					lines = append(lines, scanner.Text())
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Failed to read file:", err)
			}

			if err := file.Truncate(0); err != nil {
				fmt.Println("Failed to truncate file:", err)
			}

			if _, err := file.Seek(0, 0); err != nil {
				fmt.Println("Failed to seek file:", err)
			}

			writer := bufio.NewWriter(file)
			for _, line := range lines {
				if _, err := writer.WriteString(line + "\n"); err != nil {
					fmt.Println("Failed to write to file:", err)
				}
			}
			if err := writer.Flush(); err != nil {
				fmt.Println("Failed to flush writer:", err)
			}

			w.SetContent(hide_particular_challenge())
		})

		// ラベルと削除ボタンを横に並べる
		page.Append(widget.NewHBox(label, removeButton))
	}
	return page
}
