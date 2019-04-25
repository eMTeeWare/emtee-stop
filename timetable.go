package main

type Timetable struct {
	Station string `xml:"station,attr"`
	Stops   []struct {
		Id        string `xml:"id,attr"`
		TrainData struct {
			F           string `xml:"f,attr"`
			Type        string `xml:"t,attr"`
			Owner       string `xml:"o,attr"`
			TrainClass  string `xml:"c,attr"`
			TrainNumber int    `xml:"n,attr"`
		} `xml:"tl"`
		Arrival struct {
			PlannedTime     string `xml:"pt,attr"`
			PlannedPlatform string `xml:"pp,attr"`
			TrainLine       string `xml:"l,attr"`
			PlannedPath     string `xml:"ppth,attr"`
		} `xml:"ar"`
		Departure struct {
			PlannedTime     string `xml:"pt,attr"`
			PlannedPlatform string `xml:"pp,attr"`
			TrainLine       string `xml:"l,attr"`
			PlannedPath     string `xml:"ppth,attr"`
		} `xml:"dp"`
	} `xml:"s"`
}
