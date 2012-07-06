package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var allowanceRequestURL = "http://192.168.0.1/stlui/user/allowance_request.html"

// Regex to match the various fields, text to match looks similar to below:
// <td style="border-width:0px;">Maximum Download Bank (MB)</td><td style="border-width:0px;">1050</td>
// <td style="border-width:0px;">Allowance Remaining (MB)</td><td style="border-width:0px;">627</td>
// <td style="border-width:0px;">Allowance Remaining (%)</td><td style="border-width:0px;">68</td>
// <td style="border-width:0px;">Time Until Allowance Refill</td><td style="border-width:0px;">0:10:59:47</td>
// td style="border-width:0px;">Plan Refill Amount (MB)</td><td style="border-width:0px;">525</td>
var mbBankRegex = regexp.MustCompile("(?:<td style=\"border-width:0px;\">Maximum Download Bank \\(MB\\)</td><td style=\"border-width:0px;\">([0-9]+)</td>)")
var mbRemainingRegex = regexp.MustCompile("(?:<td style=\"border-width:0px;\">Allowance Remaining \\(MB\\)</td><td style=\"border-width:0px;\">([0-9]+)</td>)")
var pctRemainingRegex = regexp.MustCompile("(?:<td style=\"border-width:0px;\">Allowance Remaining \\(%\\)</td><td style=\"border-width:0px;\">([0-9]+)</td>)")
var timeRefillRegex = regexp.MustCompile("(?:<td style=\"border-width:0px;\">Time Until Allowance Refill</td><td style=\"border-width:0px;\">([0-9]+:[0-9]+:[0-9]+:[0-9]+)</td>)")
var mbRefillAmountRegex = regexp.MustCompile("(?:td style=\"border-width:0px;\">Plan Refill Amount \\(MB\\)</td><td style=\"border-width:0px;\">([0-9]+)</td>)")

func main() {
	res, err := http.Get(allowanceRequestURL)
	if err != nil {
		log.Fatal(err)
	}
	allowanceReq, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	mbBank := findSubmatch(mbBankRegex, allowanceReq)
	mbRemaining := findSubmatch(mbRemainingRegex, allowanceReq)
	pctRemaining := findSubmatch(pctRemainingRegex, allowanceReq)
	timeRefill := findSubmatch(timeRefillRegex, allowanceReq)
	mbRefillAmount := findSubmatch(mbRefillAmountRegex, allowanceReq)

	fmt.Printf("%s\n", allowanceReq)
	fmt.Printf("Maximum Download Bank (mb): %s\n", mbBank)
	fmt.Printf("Allowance Remaining (mb): %s\n", mbRemaining)
	fmt.Printf("Allowance Remaining (percent): %s\n", pctRemaining)
	fmt.Printf("Time until refill: %s\n", timeRefill)
	fmt.Printf("Refill Amount (mb): %s\n", mbRefillAmount)
}

func findSubmatch(re *regexp.Regexp, b []byte) (retByte []byte) {
	matches := re.FindSubmatch(b)
	if len(matches) != 2 {
		log.Fatal("No match")
	}
	retByte = matches[1]
	return
}

