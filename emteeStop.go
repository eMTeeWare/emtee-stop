package main

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var bearerToken string

func init() {
	log.SetLevel(log.DebugLevel)
	if len(os.Args) < 2 {
		fmt.Println("Please provide your authorization token as the first parameter of the application")
		os.Exit(401)
	}
	bearerToken = os.Args[1]
}

func main() {
	stationCode := "8000238"
	date := getCurrentDateForQuery()
	hour := getNextFullHourForQuery()

	var timetable Timetable
	requestDataFromDbApi(&timetable, stationCode, date, hour)

	hour = getCurrentFullHourForQuery()
	requestDataFromDbApi(&timetable, stationCode, date, hour)

	hour = getPreviousFullHourForQuery()
	requestDataFromDbApi(&timetable, stationCode, date, hour)

	for index, stop := range timetable.Stops {
		departure := stop.Departure
		trainData := stop.TrainData
		if strings.Contains(departure.PlannedPath, "Hamburg") && strings.Contains(trainData.TrainClass, "ME") {
			fmt.Printf("%02d: %s: %s von Gleis %s\t%s\n", index, formatTimeFromApiTimestamp(departure.PlannedTime), departure.TrainLine, departure.PlannedPlatform, departure.PlannedPath)
		}
	}
}

func getPreviousFullHourForQuery() string {
	return fmt.Sprintf("%02d", time.Now().Hour()-1)
}

func getCurrentFullHourForQuery() string {
	return fmt.Sprintf("%02d", time.Now().Hour())
}

func requestDataFromDbApi(timetable *Timetable, stationCode string, date string, hour string) {
	apiUrl := "https://api.deutschebahn.com/timetables/v1/plan/" + stationCode + "/" + date + "/" + hour
	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set("Authorization", "Bearer "+bearerToken)
	client := &http.Client{}
	log.Debug("Requesting " + request.URL.String())
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else if response.StatusCode < 200 || response.StatusCode > 299 {
		fmt.Printf("Server responded with %s\n", response.Status)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		responseString := normalizeXml(data)

		err = xml.Unmarshal([]byte(responseString), &timetable)
		if err != nil {
			fmt.Printf("XML parsing error: %s\n", err)
		} else {
			log.Debug(timetable)
		}
	}
}

func formatTimeFromApiTimestamp(timestamp string) string {
	return string(timestamp[6:8]) + ":" + string(timestamp[8:]) + " Uhr"
}

func getNextFullHourForQuery() string {
	return fmt.Sprintf("%02d", time.Now().Hour()+1)
}

func getCurrentDateForQuery() string {
	t := time.Now()
	return strconv.Itoa(t.Year())[2:] + fmt.Sprintf("%02d", int(t.Month())) + fmt.Sprintf("%02d", t.Day())
}

func normalizeXml(input []byte) string {
	responseString := string(input)
	responseString = strings.Replace(responseString, "><", ">\n<", -1)
	responseString = strings.Replace(responseString, "&#252;", "ü", -1)
	responseString = strings.Replace(responseString, "&#246;", "ö", -1)
	responseString = strings.Replace(responseString, "&#223;", "ß", -1)
	return responseString
}
