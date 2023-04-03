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

func readuserInfo() (string, string) {
	file, err := os.Open("userInfo.csv") 
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll() // csvを一度に全て読み込む
	if err != nil {
		log.Fatal(err)
	}
	studentId = rows[0][0]
	password = rows[0][1]

	// [][]stringなのでループする

	return password, studentId
}
