package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	tm, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	format := tm.Format("15:04:05.00 Monday 02.01.2006")
	fmt.Println(format)
}
