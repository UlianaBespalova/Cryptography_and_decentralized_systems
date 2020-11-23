package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/biter777/countries"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

var nodeUrl = "https://node1.zbeos.com"

type Data struct {
	producers []Produser
	total_producer_vote_weight string
}
type Produser struct {
	owner string
	percent float64
	producer_key string
	url string
	unpaid_blocks int
	last_claim_time time.Time
	location string
}


func httpGetProducers (url string, limit int) (map[string] interface{}, error)  { //получить от ноды данные о продюсерах
	postMessage := map[string] interface{}{
		"limit": limit,
		"lower-bound": limit,
		"json": true,
	}
	bytesRequest, err := json.Marshal(postMessage)
	if err != nil {
		return nil, err
	}

	responseUrl := url+"/v1/chain/get_producers"
	resp, err := http.Post(responseUrl, "application/json", bytes.NewBuffer(bytesRequest))
	if err != nil {
		return nil, err
	}

	var result map[string] interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}


func PercentageOfVotes (totalStr string, partStr string) (float64, error)  { //посчитать процент голосов
	total, err := strconv.ParseFloat(totalStr, 10)
	if err != nil {
		return 0, err
	}
	part, err := strconv.ParseFloat(partStr, 10)
	if err != nil {
		return 0, err
	}
	res := part/total * 100
	return res, nil
}


func GetProducers (limit int) (*Data, error) {

	res, err := httpGetProducers(nodeUrl, limit)
	if err != nil {
		return nil, err
	}
	if res==nil {
		return nil, errors.New("Failed to get data about producers\n")
	}

	var data Data //вытаскиваем нужные данные
	data.total_producer_vote_weight = res["total_producer_vote_weight"].(string)

	for _, producer_i := range res["rows"].([]interface{}) {
		producer := producer_i.(map[string]interface{})

		var newProducer Produser
		newProducer.owner = producer["owner"].(string)
		newProducer.url = producer["url"].(string)
		newProducer.producer_key = producer["producer_key"].(string)
		newProducer.unpaid_blocks = int(producer["unpaid_blocks"].(float64))
		newProducer.location = countries.ByNumeric(int((producer["location"].(float64)))).String()
		newProducer.last_claim_time, _ = time.Parse("2006-01-02T15:04:05", producer["last_claim_time"].(string))
		//--------------
		percent, err := PercentageOfVotes(data.total_producer_vote_weight, producer["total_votes"].(string))
		if err != nil {
			newProducer.percent = -1
		} else {
			newProducer.percent = percent
		}

		data.producers = append(data.producers, newProducer)
	}
	return &data, nil
}

func TableOutput (data *Data) {
	fmt.Printf("%5s%-15s%-32s%-20s%-8s%-21s%-10s\n",
		"", "ACCOUNT", "URL", "LOCATION", "PCT", "LAST CLAIM TIME", "PRODUCER KEY")
	for i, producer := range data.producers {
		fmt.Printf("%-4d %-14s %-31s %-18s %.2f%% %21s   %-10s\n",
			i+1, producer.owner,
			producer.url,
			producer.location,
			producer.percent,
			producer.last_claim_time.Format("15:04:05 Jan-02-06"),
			producer.producer_key)
	}
}


func main() {
	data, err:= GetProducers(21) //выводим список основных блок-продюсеров
	if err != nil {
		fmt.Println(err)
		return
	}
	TableOutput(data)
}
