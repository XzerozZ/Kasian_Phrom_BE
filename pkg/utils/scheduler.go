package utils

import (
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

func callTransactionAPI() {
	_, err := http.Post("http://localhost:5000/transaction/all", "application/json", nil)
	if err != nil {
		log.Println("Failed to call transaction API:", err)
	} else {
		log.Println("Transaction API called successfully")
	}
}

func StartScheduler() {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("0 0 1 * *", callTransactionAPI)
	if err != nil {
		log.Fatal("Failed to schedule job:", err)
	}

	c.Start()
	log.Println("Cron job started...")

	go func() {
		select {}
	}()
}
