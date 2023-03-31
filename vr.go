package main

import (
	"fmt"
	"time"
)

func main() {
	daytime := "26.11.2022 18:20:00"
	t, err := time.Parse("2.1.2006 15:04:05", daytime)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	vrem := int(t.Unix())
	//vrem := 1669486810
	fmt.Println("Hello, 世界", t, vrem)
	vrem2 := int64(vrem) - 6*3600
	ut := time.Unix(vrem2, 0)
	fmt.Println("\n", ut.Format("02 Jan 2006 15:04:05"))
}
