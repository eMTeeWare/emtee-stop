package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	client := &http.Client{}
	stationCode := "8000238"
	date := getCurrentDateForQuery()
	hour := getNextFullHourForQuery()
	apiUrl := "https://api.deutschebahn.com/timetables/v1/plan/" + stationCode + "/" + date + "/" + hour
	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)

	if len(os.Args) < 2 {
		fmt.Println("Please provide your authorization token as the first parameter of the application")
		os.Exit(401)
	}
	bearerToken := os.Args[1]

	request.Header.Set("Authorization", "Bearer "+bearerToken)

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else if response.StatusCode < 200 || response.StatusCode > 299 {
		fmt.Printf("Server responded with %s\n", response.Status)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		responseString := normalizeXml(data)

		type Timetable struct {
			Station string `xml:"station,attr"`
			Stops   []struct {
				Id        string `xml:"id,attr"`
				TrainData []struct {
					F           string `xml:"f,attr"`
					Type        string `xml:"t,attr"`
					Owner       string `xml:"o,attr"`
					TrainClass  string `xml:"c,attr"`
					TrainNumber int    `xml:"n,attr"`
				} `xml:"tl"`
				Arrival []struct {
					PlannedTime     string `xml:"pt,attr"`
					PlannedPlatform string `xml:"pp,attr"`
					TrainLine       string `xml:"l,attr"`
					PlannedPath     string `xml:"ppth,attr"`
				} `xml:"ar"`
				Departure []struct {
					PlannedTime     string `xml:"pt,attr"`
					PlannedPlatform string `xml:"pp,attr"`
					TrainLine       string `xml:"l,attr"`
					PlannedPath     string `xml:"ppth,attr"`
				} `xml:"dp"`
			} `xml:"s"`
		}

		var timetable Timetable
		err = xml.Unmarshal([]byte(responseString), &timetable)
		if err != nil {
			fmt.Printf("XML parsing error: %s\n", err)
		} else {
			fmt.Println(timetable)
		}

	}
}

func getNextFullHourForQuery() string {
	hour := fmt.Sprintf("%02d", time.Now().Hour()+1)
	return hour
}

func getCurrentDateForQuery() string {
	t := time.Now()
	date := strconv.Itoa(t.Year())[2:] + fmt.Sprintf("%02d", int(t.Month())) + fmt.Sprintf("%02d", t.Day())
	return date
}

func normalizeXml(input []byte) string {
	responseString := string(input)
	responseString = strings.Replace(responseString, "><", ">\n<", -1)
	responseString = strings.Replace(responseString, "&#252;", "ü", -1)
	responseString = strings.Replace(responseString, "&#246;", "ö", -1)
	responseString = strings.Replace(responseString, "&#223;", "ß", -1)
	return responseString
}
