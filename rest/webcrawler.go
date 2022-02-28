package rest

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"webcrawler/urls"
	"webcrawler/utils"
)

var scan_token = make(chan struct{}, 20)

// webCrawler : handler function for webcrawler endpoint. It trigger the processing of the sitemap
func (s *Server) webCrawler(w http.ResponseWriter, r *http.Request) {

	// Common errors to return
	var badRequestError = &Error{ID: "bad_request", Status: 400, Title: "Bad request", Detail: "No based url specified."}
	var badURLError = &Error{ID: "bad_request", Status: 400, Title: "Bad request", Detail: "Provided url can't be parsed."}

	// indicate that the serve ris processing a client request
	//Ongoing = true

	// get the base URL from the client request
	var baseurl string
	if value, ok := r.URL.Query()["url"]; ok {
		baseurl = strings.TrimSpace(value[0])
	} else {
		log.Printf("No based url specified.\n")
		writeError(w, r, badRequestError)
		return
	}

	// set baseurl to corresponding root url (e.g https://www.foo.com/bar will  https://www.foo.bar) if full parameter
	// has been provided and is equal to "true"

	if value, ok := r.URL.Query()["full"]; ok {
		full := strings.ToLower(strings.TrimSpace(value[0]))
		if full == "true" {
			baseurl = utils.GetRootURL(baseurl)
		}
	}

	base, err := url.Parse(baseurl)

	if err != nil {
		log.Printf(err.Error() + "\n")
		badRequestError.Detail = err.Error()
		writeError(w, r, badURLError)
		return
	}
	baseurl = base.String()

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
			urls.Access.Lock()
			if !alreadyScanned[url.Url] {
				n++
				//alreadyScanned[url.Url] = true
				urls.Access.Unlock()
				go func(url *urls.Urlsstore) {

					urlslist <- scanUrl(url, base, alreadyScanned)
				}(url)
			} else {
				urls.Access.Unlock()
			}
		}
	}

	// for debugging purpose we print the sitemap before returning
	// it can be commented/uncommented without impact on sitemap computation
	// urls.PrintNode(&root, 0)

	// encode the root node and seni it in to the client
	encodeJSONResponse(w, r, root)

}

// scanUrl : this method with until it can acquire a lock (currently the number of token is fixed to 20
// but it could be a parameter provided at server start. When a token is acquired, it computes the sitemap of
// passed url myurl (takining into account the initial base url baseurl). Then the list of URls us passed in the urlslist
// to be processed in turn
func scanUrl(myurl *urls.Urlsstore, baseurl *url.URL, alreadyScanned map[string]bool) []*urls.Urlsstore {
	scan_token <- struct{}{} // we wait to acquire a token

	// debug purpose only log, should be removed if we don't need it
	//fmt.Printf("DEBUG: ScanUrl  %s\n", myurl.Url)

	// we get the list of links of the considered url (they will be added on the urlslist  channel)
	list, err := urls.Sitemap(myurl, baseurl, alreadyScanned)
	<-scan_token // no it's time to release our token

	if err != nil {
		log.Print(err)
	}
	return list
}
