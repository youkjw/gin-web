package main

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		//models.CleanAllTag()

		p := "哈哈哈"
		fmt.Println(p, &p)

		p = "哈哈哈122"
		fmt.Println(p, &p)
	})

	//c.AddFunc("* * * * * *", func() {
	//	log.Println("Run models.CleanAllArticle...")
	//	//models.CleanAllArticle()
	//})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
			case <-t1.C:
				t1.Reset(time.Second * 10)
		}
	}
}