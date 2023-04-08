package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

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

		box4 := widget.NewVBox(

			widget.NewButton("特定の課題を非表示にする", func() {
				w.SetContent(hide_particular_challenge())
			}))
		page.Append(widget.NewGroup("ID・パスワード入力", box2))

		page.Append(widget.NewGroup("URL入力", box1))

		if s := exportTime(); s != "" {
			box3 := widget.NewVBox(widget.NewLabel("URL最終更新時間：" + s + "\n" + "60日でURLの期限が切れます。"))
			page.Append(box3)
		} else {
			box3 := widget.NewVBox(widget.NewLabel("URL未登録または予測していないエラーが発生しています"))
			page.Append(box3)
		}
		page.Append(widget.NewGroup("その他の設定", box4))
		url_h, err := url.Parse("https://hoppii.hosei.ac.jp/portal")
		url_m, err := url.Parse("https://cms.cis.k.hosei.ac.jp/")

		if err != nil {
			fmt.Println("URL is not found")
		}
		box5 := widget.NewVBox(widget.NewHyperlink("hoppii", url_h), widget.NewHyperlink("Moodle", url_m))

		page.Append(widget.NewGroup("関連リンク", box5))

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
			assignmentBox.Append(widget.NewLabel(fmt.Sprintf("%04d/%02d/%02d %02d:%02d \n・%s", v.time.Year, v.time.Month, v.time.Day, v.time.Hour, v.time.Min, v.name)))
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
func regist(inputID *widget.Entry, inputPass *widget.Entry) {
	fmt.Println(inputID.Text, inputPass.Text)
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

	record := []string{inputID.Text, inputPass.Text}
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
