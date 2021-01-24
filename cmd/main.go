package main

import (
	"fmt"
	"time"
	_ "time/tzdata"
	"workshop"
	_ "workshop/docs"
)

func main() {
	fmt.Print("Local time zone ")
	fmt.Println(time.Now().Zone())
	fmt.Println(time.Now())
	workshop.StartServer()
}
