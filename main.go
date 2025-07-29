package main

import (
	"fmt"

	"github.com/R0Xps/gatorcli/internal/config"
)

func main() {
	conf, err := config.Read()
	fmt.Println(conf)
	fmt.Println(err)
	conf.SetUser("yahya")
	conf, err = config.Read()
	fmt.Println(conf)
	fmt.Println(err)
}
