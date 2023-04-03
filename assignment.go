package main

import (
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

	// // Announcement Entity の User Action を取得
	// auo, err := hoppiiApi.Announcement.User(&announcement.AnnouncementUserInput{
	// 	UserCookie: user.Cookie,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// // 結果を出力
	// fmt.Println(auo.AnnouncementUserResult.AnnouncementCollection[0].Title)

	// Assignment Entity の My Action を取得
	ao, err := hoppiiApi.Assignment.My(&assignment.AssignmentMyInput{
		UserCookie: user.Cookie,
	})
	if err != nil {
		return nil
	}
	// 結果を出力
	for _, v := range ao.AssignmentMyResult.AssignmentCollection {
		var info classInfo

		info.name = v.Title
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

		if time.Now().Month() > info.time.Month {
			continue
		}

		classInfos = append(classInfos, info)
	}
	return classInfos

}
