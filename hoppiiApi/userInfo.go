package main

import (
	"encoding/csv"
	"log"
	"os"
)

var (
	password  string
	studentId string
)

func readuserInfo() (string,string) {
	file, err := os.Open("../userInfo.csv") // 先ほど入手した郵便番号データをos.Openで開く
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll() // csvを一度に全て読み込む
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range rows[0] {
		password += v
	}
	for _, v := range rows[1] {
		studentId += v
	}
	// [][]stringなのでループする

	return password, studentId
}
