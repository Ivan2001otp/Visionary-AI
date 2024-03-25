package main

import (
	"fmt"
	"log"

	// "strconv"

	"sync"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
	scrap "github.com/Ivan2001otp/Visionary-AI/Service/Scrap"
	// config "github.com/Ivan2001otp/Visionary-AI/Service/Database"
)

func storePersistently(items []model.Product) {
	for i, item := range items {
		result, err := item.SaveToMongo()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted ", i+1, "->", result)
	}
}

func InitProcess() {

	//creating Channels that to signal completition of each goroutine..
	doneChannel := make(chan struct{})
	model.DeleteAllFromMongo()
	//wait group.

	var wg sync.WaitGroup //-> these helps to lock the Critical section during concurrent operation.

	wg.Add(2)
	startTime := time.Now()

	//use of go-routines that are useful to carry out async tasks.

	// This goroutine service fetches the television info
	go func() {

		defer wg.Done()
		result, err := scrap.TelevisionScraper()
		if err != nil {

			log.Fatal("Error catched in driver function while scraping television-data->", err)
		}

		//save to mongodb->...

		fmt.Println("the Tele-len is ", len(result))
		doneChannel <- struct{}{}
	}()

	<-doneChannel //emit the signal signifying that the first goroutine has done the job..

	//this gorouting scrapes the smartphone information.
	go func() {
		defer wg.Done()

		result, err := scrap.SmartPhoneScrapper()

		if err != nil {
			log.Fatal("Error catched in driver function while scraping television-data->", err)
		}

		//save to mongodb...
		fmt.Println("The response len is ", len(result))
		doneChannel <- struct{}{}
	}()

	<-doneChannel //emit the signal signifying that the second goroutine has done the job..

	/*
		//this go routine fetches watches info.

		go func() {

			defer wg.Done()
			result, err := scrap.WatchScraper()

			if err != nil {
				log.Fatal("Error catched in driver function while scraping smartphone-data->", err)
			}

			fmt.Println("the response has ", len(result))
			doneChannel <- struct{}{}

		}()
	*/

	// scraps the info about laptops
	/*
		go func() {
			defer wg.Done()
			result, err := scrap.LapTopScraper()

			if err != nil {
				log.Fatal("Error catched in driver function while scraping Laptop-data->", err)
			}

			//save to mongo...

			fmt.Println("The response has ->", len(result))
			doneChannel <- struct{}{}
		}()
	*/

	wg.Wait()

	endTime := time.Now()
	fmt.Println("The time taken is ", endTime.Sub(startTime))

}

func main() {

	//invoke the service part once in 10 minutes..

	ticker := time.NewTicker(10 * time.Minute)

	defer ticker.Stop()

	/*

		Note : Do not run this InitProcess function,it is computationally intensive task..
		..Kindly ask permission before invoking the InitProcess..The bugs are still in progress to be fixed.
		Thank you.


			for range ticker.C{
			InitProcess();//->the function which is necessary to invoke the possible scraper service...
			}
	*/

	select {}

}
