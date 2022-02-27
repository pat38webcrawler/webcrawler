package rest

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"webcrawler/urls"
)

var Ongoing = false

var scan_token = make(chan struct{}, 20)

// webCrawler : handler function for webcrawler endpoint. It trigger the processing of the sitemap
func (s *Server) webCrawler(w http.ResponseWriter, r *http.Request) {

	// Common errors to return
	var badRequestError = &Error{ID: "bad_request", Status: 400, Title: "Bad request", Detail: "No based url specified."}
	var badURLError = &Error{ID: "bad_request", Status: 400, Title: "Bad request", Detail: "Provided url can't be parsed."}
	var ServiceUnavailableError = &Error{ID: "bad_request", Status: 503, Title: "Service Unvailable", Detail: "The web crawler is already running."}

	// check if a request is already process by the server and return if it is the case
	// it's a design choice for the ake of simplicity, it could be improve to run request in //
	if Ongoing == true {
		writeError(w, r, ServiceUnavailableError)
		return
	}

	// indicate that the serve ris processing a client request
	Ongoing = true

	// get the base URL from the client request
	var baseurl string
	if value, ok := r.URL.Query()["url"]; ok {
		baseurl = strings.TrimSpace(value[0])
	} else {
		log.Printf("No based url specified.\n")
		writeError(w, r, badRequestError)
		Ongoing = false
		return
	}

	base, err := url.Parse(baseurl)
	if err != nil {
		log.Printf(err.Error() + "\n")
		badRequestError.Detail = err.Error()
		writeError(w, r, badURLError)
		Ongoing = false
		return
	}

	//initialize the root Node of our sitemap
	// this node will be return to the client at the end of the processing
	var root urls.Node
	root.Url = baseurl

	// channel : contain the list of url to crawl.
	urlslist := make(chan []*urls.Urlsstore)

	var n int = 1
	iurl := &urls.Urlsstore{Url: baseurl, Node: &root, Skip: false}

	// we initialize the channel array with the base url provided by the client
	go func() { urlslist <- []*urls.Urlsstore{iurl} }()

	// this map will keep already seen urls in order to avoid cyclic processing
	alreadyScanned := make(map[string]bool)

	for ; n > 0; n-- {

		list := <-urlslist

		for _, url := range list {
			if !alreadyScanned[url.Url] {
				n++
				alreadyScanned[url.Url] = true
				go func(url *urls.Urlsstore) {

					urlslist <- scanUrl(url, base)
				}(url)
			}
		}
	}

	// for debugging purpose we print the sitemap before returning
	// it can be commented
	urls.PrintNode(&root, 0)

	// encode the root node and seni it in to the client
	encodeJSONResponse(w, r, root)

	// we reinitialize the VisitedUrls map for the nrxt request
	urls.VisitedURLs = make(map[string]bool)
	// and we inidctes that the server is ready to handle a new request
	Ongoing = false
}

// scanUrl : this method with until it can acquire a lock (currently the number of token is fixed to 20
// but it could be a parameter provided at server start. When a token is acquired, it computes the sitemap of
// passed url myurl (takining into account the initial base url baseurl). Then the list of URls us passed in the urlslist
// to be processed in turn
func scanUrl(myurl *urls.Urlsstore, baseurl *url.URL) []*urls.Urlsstore {
	scan_token <- struct{}{} // we wait to acquire a token

	// debug purpose only log, should be removed if we don't need it
	fmt.Printf("DEBUG: ScanUrl  %s\n", myurl.Url)

	// we get the list of links of the considered url (they will be added on the urlslist  channel)
	list, err := urls.Sitemap(myurl, baseurl)
	<-scan_token // no it's time to release our token

	if err != nil {
		log.Print(err)
	}
	return list
}
