package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

func main() {
	log.Println("Starting...")

	var wg sync.WaitGroup
	wg.Add(1)

	c := cron.New(cron.WithSeconds())
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag1...")
		//models.CleanAllTag()
	})

	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag2...")
		//models.CleanAllTag()
	})

	//c.AddFunc("* * * * * *", func() {
	//	log.Println("Run models.CleanAllArticle...")
	//	//models.CleanAllArticle()
	//})

	c.Start()
	defer c.Stop()

	//t1 := time.NewTimer(time.Second * 2)
	//for {
	//	select {
	//	case <-t1.C:
	//		t1.Reset(time.Second * 2)
	//	}
	//}

	wg.Wait()
}