package main

import (
	"github.com/RichardKnop/voucher/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		panic(err)
	}
}
