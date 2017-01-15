package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"io/ioutil"
	"strings"
	"encoding/xml"
	"time"
	"os"
	"strconv"
)

type AllSms struct {
	XMLName   xml.Name `xml:"allsms"`
	Count     string   `xml:"count,attr"`
	Sms       []Sms
}

type Sms struct {
	XMLName   xml.Name `xml:"sms"`
	Address   string   `xml:"address,attr"`
	Time      string   `xml:"time,attr"`
	Date      string   `xml:"date,attr"`
	Type      string   `xml:"type,attr"`
	Body      string   `xml:"body,attr"`
	Read      string   `xml:"read,attr"`
}

func parseDate(dateStr string) (time.Time) {
	form := "01/02/2006 15:04:00"
	parsedTime, err := time.Parse(form, dateStr)
	if err != nil {
		form := "01/02/2006 15:04"
		parsedTime, err := time.Parse(form, dateStr)

		if err != nil {
			log.Fatal(err)
		}

		return parsedTime
	}

	return parsedTime
}

func parseDirection(readStatus string) (string)  {
	if (readStatus == "R") {
		return "1"
	} else {
		return "2"
	}
}

func main() {
	in, err := ioutil.ReadFile("export-avast-backup-sms.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(in)))

	var sms []Sms
	for  {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		address := record[0]
		direction := parseDirection(record[1])
		datetime:= parseDate(record[2])
		body := record[3]

		newSMS := Sms{
			Address: address,
			Time: datetime.Format("2 Jan 2006 15:37:15"),
			Date: strconv.FormatInt(datetime.Unix(), 10),
			Type: direction,
			Body: body,
			Read: "1",
		}
		sms = append(sms, newSMS)

		fmt.Println(record[0])
		fmt.Println(record[1])
		fmt.Println(record[2])
		fmt.Println(record[3])
		fmt.Println("----- Another ------")
	}

	AllSms := AllSms{Sms: sms, Count: strconv.Itoa(len(sms))}

	f, err := os.Create("out.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := xml.NewEncoder(f)
	enc.EncodeElement(`<?xml version="1.0" encoding="UTF-8"?>` + "\n", xml.StartElement{})
	if err := enc.Encode(AllSms); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
