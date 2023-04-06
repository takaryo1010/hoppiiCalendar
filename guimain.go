package main

import (
	"bufio"
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

func registerID(page *widget.Box) {
	if currentPage == 2 {
		inputID := widget.NewEntry()
		inputpass := widget.NewPasswordEntry()
		inputURL := widget.NewEntry()

		box1 := widget.NewVBox(
			widget.NewForm(
				widget.NewFormItem("URL", inputURL),
			),
			widget.NewButton("登録", func() {
				registURL(inputURL.Text)
				w.SetContent(NewPage(currentPage))
			}))

		box2 := widget.NewVBox(
			widget.NewForm(
				widget.NewFormItem("ID", inputID),
				widget.NewFormItem("PASS", inputpass),
			),
			widget.NewButton("登録", func() {
				regist(inputID, inputpass)
				w.SetContent(NewPage(currentPage))
			}))
		// TODO
		box4 := widget.NewVBox(

			widget.NewButton("特定の課題を非表示にする", func() {
				w.SetContent(hide_particular_challenge())
			}))
		page.Append(widget.NewGroup("ID・パスワード入力",box2))

		page.Append(widget.NewGroup("URL入力",box1))

		if s := exportTime(); s != "" {
			box3 := widget.NewVBox(widget.NewLabel("URL最終更新時間：" + s + "\n" + "60日でURLの期限が切れます。"))
			page.Append(box3)
		} else {
			box3 := widget.NewVBox(widget.NewLabel("URL未登録または予測していないエラーが発生しています"))
			page.Append(box3)
		}
		page.Append(widget.NewGroup("その他の設定",box4))
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
func regist(inputID *widget.Entry, inputpass *widget.Entry) {
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
}

func registURL(s string) {
	file, err := os.Create("URL.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// CSVライターを作成する
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSVにデータを書き込む
	dateStr := time.Now().Format("2006-01-02 15:04:05")
	writer.Write([]string{s, dateStr})
}
func exportTime() string {
	_, err := os.Stat("URL.csv")
	if os.IsNotExist(err) {
		fmt.Println("URL.csv file does not exist.")
		return ""
	}

	// URL.csvファイルを読み込む
	file, err := os.Open("URL.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	// CSVパーサーを作成する
	csvReader := csv.NewReader(file)

	// URL.csvファイル内の時間を出力する
	for {
		record, err := csvReader.Read()
		if err != nil {
			fmt.Println("Error reading file:", err)
			return ""
		}

		// record[0]が時間の文字列であることを前提としている
		t, err := time.Parse("2006-01-02 15:04:05", record[1])
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}

		return t.Format("2006-01-02 15:04:05")
	}

}
