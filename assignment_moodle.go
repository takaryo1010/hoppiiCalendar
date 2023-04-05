package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func ics() []classInfo {
	// Calendar.iscファイルを読み込む
	if !getIcs() {
		var a []classInfo
		a = append(a, classInfo{"moodleは登録されていません", timeInfo{time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second()}, true})
		return a
	}

	file, err := ioutil.ReadFile("calendar.ics")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content := string(file)

	// SUMMARYに「終了」が含まれるイベントを検索し、DTENDとSUMMARYを格納する
	re := regexp.MustCompile("(?s)BEGIN:VEVENT.*?SUMMARY:(.*?)\\n.*?DTEND:(\\d{8}T\\d{6}Z).*?END:VEVENT")
	matches := re.FindAllStringSubmatch(content, -1)
	events := []classInfo{}
	for _, match := range matches {
		var info classInfo
		summary := match[1]
		if strings.Contains(summary, "終了") {
			endTime, err := time.Parse("20060102T150405Z", match[2])
			if err != nil {
				fmt.Println(err)
				continue
			}
			// UTCからJSTに変換する
			jst := time.FixedZone("Asia/Tokyo", 9*60*60)
			endTime = endTime.In(jst)
			var dateTime timeInfo
			dateTime.Year, dateTime.Month, dateTime.Day = endTime.Date()
			dateTime.Hour, dateTime.Min, dateTime.Sec = endTime.Clock()
			info.name = "M:" + summary
			info.status = true
			info.time = dateTime
			if endTime.After(time.Now()) {
				events = append(events, info)
			}

		}
	}
	return events
}

func getIcs() bool {
	password, username := readuserInfo()
	url := readuserURL()

	fileName := "calendar.ics"

	// BASIC認証情報を含むHTTPクライアントを作成
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.SetBasicAuth(username, password)

	// URLからファイルをダウンロード
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// ファイルを作成
	file, err := os.Create(fileName)
	if err != nil {
		return false
	}
	defer file.Close()

	// ファイルを保存
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return false
	}
	return true
}
