package main

import (
	"bufio"
	"fmt"
	"os"

	"l2.9/internal/unpack"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var str string
	scanner.Scan()
	str = scanner.Text()
	upd, err := unpack.String(str)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	fmt.Println(upd)
}
