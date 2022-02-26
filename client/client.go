package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"webcrawler/urls"
)

// client code

func main() {
	// we check that there is exactly 2 parameters passed to the client
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <host_or_ip> <https://myserver.com> ", os.Args[0])
	}
	host := os.Args[1]
	urltoCrawl := os.Args[2]

	// get an http.Client
	client, err := client(host)
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
