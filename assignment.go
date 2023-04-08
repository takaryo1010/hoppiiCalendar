package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	hoppii_api "github.com/isso-719/hoppii-api"
	"github.com/isso-719/hoppii-api/pkg/assignment"
	"github.com/isso-719/hoppii-api/pkg/user"
)

type timeInfo struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
	Min   int
	Sec   int
}

type classInfo struct {
	name   string
	time   timeInfo
	status bool
}

func assigntmentInfo() []classInfo {
	var (
		classInfos []classInfo
	)
	password, studentId := readuserInfo()
	hoppiiApi := hoppii_api.NewHoppiiApi()

	// ユーザクライアントの作成 (これを使用して Hoppii にアクセスする)
	user, err := hoppiiApi.User.CreateUser(&user.UserInput{
		StudentId: studentId,
		Password:  password,
	})
	if err != nil {

		return nil
	}
	ao, err := hoppiiApi.Assignment.My(&assignment.AssignmentMyInput{
		UserCookie: user.Cookie,
	})
	if err != nil {
		return nil
	}
	// 結果を出力
	for _, v := range ao.AssignmentMyResult.AssignmentCollection {
		var info classInfo

		info.name = "H:" + v.Title
		if v.Status != "CLOSED" {
			info.status = true
		}
		// 2022-06-16 01:40:00 +0000 UTC
		utcTime := v.CloseTimeString
		t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", utcTime.String())
		if err != nil {
			panic(err)
		}

		jst := time.FixedZone("Asia/Tokyo", 9*60*60) // JSTはUTC+9
		localTime := t.In(jst)
		var dateTime timeInfo
		dateTime.Year, dateTime.Month, dateTime.Day = localTime.Date()
		dateTime.Hour, dateTime.Min, dateTime.Sec = localTime.Clock()
		info.time = dateTime

		// 日時が現在時刻よりも新しいものだけ classInfos に追加する
		if localTime.After(time.Now()) {
			classInfos = append(classInfos, info)
		}
	}

	classInfos = append(classInfos, ics()...)
	classInfos, err = removeElementsFromSliceByCSV("remove.csv", classInfos)
	if err != nil {
		log.Println("Could not process remove.csv")
	}
	classInfos = sortClassesByTime(classInfos)
	return classInfos

}

func removeElementsFromSliceByCSV(filename string, slice []classInfo) ([]classInfo, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		for i := 0; i < len(slice); i++ {
			if strings.Contains(slice[i].name, record[0]) {
				slice = append(slice[:i], slice[i+1:]...)
				i--
			}
		}
	}

	return slice, nil
}

func sortClassesByTime(classes []classInfo) []classInfo {
	sort.Slice(classes, func(i, j int) bool {
		return time.Date(classes[i].time.Year, classes[i].time.Month, classes[i].time.Day, classes[i].time.Hour, classes[i].time.Min, classes[i].time.Sec, 0, time.Local).Before(
			time.Date(classes[j].time.Year, classes[j].time.Month, classes[j].time.Day, classes[j].time.Hour, classes[j].time.Min, classes[j].time.Sec, 0, time.Local))
	})
	return classes
}
