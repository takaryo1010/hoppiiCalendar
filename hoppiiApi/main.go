package main

import (
	"fmt"
	"github.com/isso-719/hoppii-api"
	"github.com/isso-719/hoppii-api/pkg/announcement"
	"github.com/isso-719/hoppii-api/pkg/assignment"
	"github.com/isso-719/hoppii-api/pkg/user"
)

func main() {
	password,studentId:=readuserInfo()
	hoppiiApi := hoppii_api.NewHoppiiApi()
	
	// ユーザクライアントの作成 (これを使用して Hoppii にアクセスする)
	user, err := hoppiiApi.User.CreateUser(&user.UserInput{
		StudentId: studentId,
		Password:  password,
	})
	if err != nil {
		panic(err)
	}

	// Announcement Entity の User Action を取得
	auo, err := hoppiiApi.Announcement.User(&announcement.AnnouncementUserInput{
		UserCookie: user.Cookie,
	})
	if err != nil {
		panic(err)
	}
	// 結果を出力
	fmt.Println(auo.AnnouncementUserResult.AnnouncementCollection[0].Title)

	// Assignment Entity の My Action を取得
	ao, err := hoppiiApi.Assignment.My(&assignment.AssignmentMyInput{
		UserCookie: user.Cookie,
	})
	if err != nil {
		panic(err)
	}
	// 結果を出力
	fmt.Println(ao.AssignmentMyResult.AssignmentCollection[0].Title)
}
