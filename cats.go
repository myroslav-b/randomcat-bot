package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const cAddress = "https://api.thecatapi.com/v1/images/search"

//TCat contains information about cats imported from thecatapi.com
type TCat struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func getCat(client *http.Client) (TCat, error) {
	//var cat TCat
	cats := make([]TCat, 1)
	req, err := http.NewRequest("GET", cAddress, nil)
	if err != nil {
		return cats[0], errors.Wrap(err, "Error (getCat)")
	}
	//client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return cats[0], errors.Wrap(err, "Error (getCat)")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return cats[0], errors.Wrap(errors.New(resp.Status), "Error (getCat)")
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return cats[0], errors.Wrap(err, "Error reading response body (getCat)")
	}
	err = json.Unmarshal(body, &cats)
	if err != nil {
		return cats[0], errors.Wrap(err, "Error unmarshaling response body (getCat)")
	}
	return cats[0], nil
}

func catGenerator(chanCat chan string) {
	// Need at least one successful cat
	var reserveCat TCat
	client := &http.Client{Timeout: 10 * time.Second}
	for {
		switch {
		case len(chanCat) < cap(chanCat):
			cat, err := getCat(client)
			if err != nil {
				log.Println(err)
				chanCat <- reserveCat.URL
			} else {
				chanCat <- cat.URL
				reserveCat = cat
			}
		default:
			time.Sleep(time.Second)
		}
	}
}
