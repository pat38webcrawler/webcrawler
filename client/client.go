package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"webcrawler/urls"
	"webcrawler/utils"
)

// client code

func main() {
	// we check that there is exactly 3 parameters (host or ip, port and base url passed to the client
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <host_or_ip> <port> <https://myserver.com> ", os.Args[0])
	}
	host := os.Args[1]
	port := os.Args[2]

	// run some checks on port and based url
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Specified port %s is not an int", port)
	}
	urltoCrawl := os.Args[3]
	if !utils.IsCorrrectURL(urltoCrawl) {
		log.Fatalf("Specified base url  '%s' seems not correct (should be defined as 'http[s]://mysite'", port)
	}

	// get an http.Client
	client, err := client(host, port)
	if err != nil {
		log.Fatalln("An error occurred" + err.Error())
	}

	// initialize a request calling the webcrawler endpoint
	request, err := client.NewRequest("GET", "/webcrawler", nil)
	if err != nil {
		log.Fatalln("An error occurred" + err.Error())
	}
	// Initialize the header
	request.Header.Add("Accept", "application/json")

	//Initialization of the parameter used to pass the base d nto initiate the web crawling
	params := url.Values{}
	params.Set("url", urltoCrawl)
	request.URL.RawQuery = params.Encode()

	// send the actual request
	response, err := client.Do(request)
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		log.Fatalln("An error occurred" + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected error code :  %d", response.StatusCode)
	}
	// read the response and deserialized it in the appropriate struct
	sitemap := new(urls.Node)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("An error occurred" + err.Error())
	}

	err = json.Unmarshal(body, &sitemap)
	if err != nil {
		log.Fatalln("An error occurred" + err.Error())
	}

	// we can now display the returned sitemap
	urls.PrintNode(sitemap, 0)
}
