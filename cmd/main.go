package main

import (
	"log"

	"github.com/shizuokago/blog"
	"github.com/shizuokago/blog/config"
)

func main() {
	err := blog.Start(config.AppEnginePort())
	if err != nil {
		log.Println("blog start:", err)
	}
}
