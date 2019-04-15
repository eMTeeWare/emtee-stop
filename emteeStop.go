package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	client := &http.Client{}
	stationCode := "8000238"
	date := "190415"
	hour := "20"
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
				Id string `xml:"id,attr"`
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

func normalizeXml(input []byte) string {
	responseString := string(input)
	responseString = strings.Replace(responseString, "><", ">\n<", -1)
	responseString = strings.Replace(responseString, "&#252;", "ü", -1)
	responseString = strings.Replace(responseString, "&#246;", "ö", -1)
	responseString = strings.Replace(responseString, "&#223;", "ß", -1)
	return responseString
}