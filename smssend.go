/*
# smssend - smssend is a program to send SMS messages from the commandline.

# Copyright © 2009-2012 by Denis Khabarov aka 'Saymon21'
# E-Mail: saymon at hub21 dot ru (saymon@hub21.ru)

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License version 3
# as published by the Free Software Foundation.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>
*/
package main

import (
	"os"
	"fmt"
	"flag"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"bufio"
)

func main() {
	var text string
	var ApiID = flag.String("api-id", "", "Your API-ID")
	var to = flag.String("to", "", "Phone number")
	var msgs = flag.String("message", "", "Message")
	var version = flag.Bool("v", false, "Print version information and quit")
	flag.Parse()
	if *version {
		fmt.Printf("smssend.go version 0.1\n")
		os.Exit(0)
	}
	if len(*msgs) == 0 {
		bio := bufio.NewReader(os.Stdin)	
		msg, _, _ := bio.ReadLine()
		text = string(msg)		
	} else {
		text = string(*msgs)
	}
	var sUrl = "http://sms.ru/sms/send?api_id=%s&to=%s&text=%s"	
	resp, err := http.Get(fmt.Sprintf(sUrl, *ApiID, *to, url.QueryEscape(text)))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}
	var code = strings.Split(string(body),"\n")

	switch code[0] {
		case "301":
			fmt.Printf("Invalid password or user not found\n")
			os.Exit(1)
		case "302":
			fmt.Printf("This account is not verified!\n")
			os.Exit(1)
		case "300":
			fmt.Printf("Invalid token or your IP address has been changed\n")
			os.Exit(1)
		case "230":
			fmt.Printf("Reached daily limit of messages to this number!\n")
			os.Exit(1)
		case "220":
			fmt.Printf("Service temporarily unavailable. Please try again.\n")
			os.Exit(1)
		case "200":
			fmt.Printf("Invalid API ID\n")
			os.Exit(1)
		case "201":
			fmt.Printf("Out of money\n")
			os.Exit(1)
		case "202":
			fmt.Printf("Invalid recipient\n")
			os.Exit(1)
		case "203":
			fmt.Printf("Message text not specified\n")
			os.Exit(1)
		case "204":
			fmt.Printf("Bad sender\n")
			os.Exit(1)
		case "205":
			fmt.Printf("Message too long\n")
			os.Exit(1)
		case "206":
			fmt.Printf("Day message limit reached\n")
			os.Exit(1)
		case "207":
			fmt.Printf("Can't send messages to that number\n")
			os.Exit(1)
		case "208":
			fmt.Printf("Wrong time\n")
			os.Exit(1)
		case "209":
			fmt.Printf("Recipient in blacklist")
			os.Exit(1)
		case "212":
			fmt.Printf("Invalid text encoding. Expected UTF8")
		case "100":
			os.Exit(0)
	}
}
