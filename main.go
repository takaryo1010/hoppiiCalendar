package main

import "fmt"

func main() {
	guimain()
	fmt.Println("aaa")
	classInfos := assigntmentInfo()
	if classInfos == nil {
		fmt.Println("no Infomation or network error")
	} else {
		for _, v := range classInfos {
			fmt.Println(v)
		}
	}

}
