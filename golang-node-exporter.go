package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Get_Data_and_Put_Into_File() {

	//----------------------------------------------------------------------------------
	//----------------------------------------------------------------------------------
	//               in this part, i'll get some data from an API
	//
	//               and put it into "metrics.txt" file
	//
	//              you can remove this part and add your own script
	//----------------------------------------------------------------------------------
	//----------------------------------------------------------------------------------

	for {

		// This function will get data every 5 seconds
		time.Sleep(5 * time.Second)

		// Define the structure of our API
		type data struct {
			Time struct {
				Updated    string    `json:"updated"`
				UpdatedISO time.Time `json:"updatedISO"`
				Updateduk  string    `json:"updateduk"`
			} `json:"time"`
			Disclaimer string `json:"disclaimer"`
			ChartName  string `json:"chartName"`
			Bpi        struct {
				Usd struct {
					Code        string  `json:"code"`
					Symbol      string  `json:"symbol"`
					Rate        string  `json:"rate"`
					Description string  `json:"description"`
					RateFloat   float64 `json:"rate_float"`
				} `json:"USD"`
				Gbp struct {
					Code        string  `json:"code"`
					Symbol      string  `json:"symbol"`
					Rate        string  `json:"rate"`
					Description string  `json:"description"`
					RateFloat   float64 `json:"rate_float"`
				} `json:"GBP"`
				Eur struct {
					Code        string  `json:"code"`
					Symbol      string  `json:"symbol"`
					Rate        string  `json:"rate"`
					Description string  `json:"description"`
					RateFloat   float64 `json:"rate_float"`
				} `json:"EUR"`
			} `json:"bpi"`
		}

		//----------------------------------------------------------------------------------
		// Get data from API

		url := "https://api.coindesk.com/v1/bpi/currentprice.json"
		res, getErr := http.Get(url)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		data_obj := data{}
		jsonErr := json.Unmarshal(body, &data_obj)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		//----------------------------------------------------------------------------------
		// Empty the file
		if err := os.Truncate("metrics.txt", 0); err != nil {
			log.Printf("Failed to truncate: %v", err)
		}

		//----------------------------------------------------------------------------------
		// Open the file for writing data into it

		f, err := os.OpenFile("metrics.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Can't open the file.")
		}
		defer f.Close()

		//----------------------------------------------------------------------------------
		// Write data to file

		f.WriteString("\nBitcoin_USD: ")
		price := strings.Replace(data_obj.Bpi.Usd.Rate, ",", "", -1)
		f.WriteString(price)

		f.WriteString("\nBitcoin_EUR: ")
		price = strings.Replace(data_obj.Bpi.Eur.Rate, ",", "", -1)
		f.WriteString(price)
	}
}

func main() {

	//----------------------------------------------------------------------------------
	//----------------------------------------------------------------------------------
	//
	//                I used "Get_Data_and_Put_Into_File" function to gather data
	//
	//               	and put it into a file called "metrics.txt"
	//
	//               	you can write and use your own function
	//
	//----------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------

	go Get_Data_and_Put_Into_File()

	//----------------------------------------------------------------------------------
	//----------------------------------------------------------------------------------
	//
	//               		in this part,
	//
	//               i'll run a web server, so prometheus could scrape data
	//
	//----------------------------------------------------------------------------------
	//----------------------------------------------------------------------------------

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Fatal(http.ListenAndServe(":9999", nil))

}
