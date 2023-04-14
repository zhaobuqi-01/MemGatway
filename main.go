package main

import (
	"fmt"
	"gateway/configs"
	"gateway/internal/pkg"
)

func main() {
	configs.Init()
	fmt.Println(pkg.GenSaltPassword("123456"))
}
