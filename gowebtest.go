package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	rps := flag.Int("rps", 100, "int, requests per second, default 100")
	duration := flag.Int("d",5, "int, how long to test in seconds, default 5")
	url := flag.String("url","http://topky.sk","url to test")
	timeoutcode := flag.Int("tc",408,"uint, code for timeout, default 408")
	timeout := flag.Int("t",800,"int, milliseconds to timeout, default 800")
	flag.Parse()
	log.Println("Starting")
	//executionTime := time.Duration(60) * time.Second
	iteration := *duration
	results := make([]map[int]int, 0)
	for j := 0; j < iteration; j++ {
		log.Println(j)
		start := time.Now()
		quitch := make(chan bool)
		urlsch := make(chan string)
		resultsch := make(chan int)
		go func() {
			for i := 0; i < *rps; i++ {
				//urlsch <- "http://topky.sk"
				urlsch <- *url
			}
			close(urlsch)
		}()
		go func() {
			time.Sleep(1 * time.Second)
			quitch <- true
		}()
		doc := make(map[int]int)
	Loop:
		for {
			select {
			case url :=<- urlsch:
				go func(u string){
					if len(url) > 0 {
						resultsch <- Get(u, *timeoutcode, *timeout)
					}
				}(url)
			case result :=<- resultsch:
				doc[result]++
			case <-quitch:
				close(quitch)
				log.Println("Breaking")
				break Loop
			}
		}
		end := time.Now()
		diff := end.Sub(start).Milliseconds()
		doc[-1] = int(diff)
		results = append(results, doc)
		log.Println("Exiting", j)
		//if j == 4{
		//	log.Println(j)
		//	log.Println(results)
		//	os.Exit(0)
		//}
	}
	log.Println(results)
}

func Get(urlstring string, timeoutcode int, timeout int ) int {
	header := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36"
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	request, err := http.NewRequest("GET", urlstring, nil)
	if err != nil {
		return 0
	}
	//do not forget!!!
	request.Header.Set("User-Agent", header)

	// Make request
	resp, err := client.Do(request)
	if err != nil{
		if strings.Contains(err.Error(), "Client.Timeout"){
			return timeoutcode
		}
		//log.Println(err)
		return 0
	}
	statuscode := resp.StatusCode
	return statuscode
}