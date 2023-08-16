package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	ccsv "github.com/tsak/concurrent-csv-writer"
)

type extractedJob struct {
	id        string
	title     string
	location  string
	condition string
	summary   string
}

var baseURL string = "https://www.saramin.co.kr/zf_user/search?cat_kewd=84&company_cd=0%2C1%2C2%2C3%2C4%2C5%2C6%2C7%2C9%2C10&loc_mcd=101000&panel_type=&search_optional_item=y&search_done=y&panel_count=y&preview=y&recruitSort=relation&inner_com_type=&searchword=&show_applied=&quick_apply=&except_read=&ai_head_hunting=&mainSearch=n"

func main() {
	// 시작 시간
	startTime := time.Now()

	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobswithGoRoutines(jobs)
	fmt.Println("Done, extracted", time.Since(startTime))
}

func writeJobswithGoRoutines(jobs []extractedJob) {
	csv, err := ccsv.NewCsvWriter("jobs.csv")
	checkErr(err)

	defer csv.Close()

	headers := []string{"Link", "Title", "Location", "Condition", "Summary"}
	csv.Write(headers)

	done := make(chan bool)
	for _, job := range jobs {
		go func(job extractedJob) {
			csv.Write([]string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=" + job.id, job.title, job.location, job.condition, job.summary})
			done <- true
		}(job)
	}
	for i := 0; i < len(jobs); i++ {
		<-done
	}
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&recruitSort=relation&recruitPageCount=40" + "&recruitPage=" + strconv.Itoa(page+1)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".item_recruit")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".area_job>.job_tit>a").Text())
	location := cleanString(card.Find(".area_job>.job_condition>span>a").Text())
	condition := card.Find(".area_job>.job_condition>span:not(:first-child)").Text()
	summary := cleanString(card.Find(".area_job>.job_sector>a").Text())
	c <- extractedJob{
		id:        id,
		title:     title,
		location:  location,
		condition: condition,
		summary:   summary}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
