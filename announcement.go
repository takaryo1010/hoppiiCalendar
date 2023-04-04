package main

import (
	hoppii_api "github.com/isso-719/hoppii-api"
	"github.com/isso-719/hoppii-api/pkg/announcement"
	"github.com/isso-719/hoppii-api/pkg/user"
)

type announceInfo struct {
	name    string
	content string
}

func announcementInfo() []announceInfo {
	var (
		announceInfos []announceInfo
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

	// Announcement Entity の User Action を取得
	auo, err := hoppiiApi.Announcement.User(&announcement.AnnouncementUserInput{
		UserCookie: user.Cookie,
	})
	if err != nil {
		panic(err)
	}
	// 結果を出力
	for _, v := range auo.AnnouncementUserResult.AnnouncementCollection {
		var Info announceInfo
		Info.name = v.CreatedByDisplayName
		Info.content = v.Title
		announceInfos = append(announceInfos, Info)
	}

	return announceInfos

}
