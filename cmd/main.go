package main

import (
	"log"

	"github.com/shizuokago/blog"
	"github.com/shizuokago/blog/config"
)

func main() {

	//blog start
	err := blog.Start(
		config.Port(),
		config.Project(),
		config.Datastore(),
	)

	if err != nil {
		log.Printf("blog start error: %+v\n", err)
	}
}
