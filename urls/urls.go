package urls

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"webcrawler/utils"

	"golang.org/x/net/html"
)

// this is the struct u(chained data) that is used to store the whole sitemap
// Url is the corresponding url and Children is an array of Node which represent
// the current Node's children

type Node struct {
	Url      string
	Children []*Node
}

// this is the struct used to described a URL that wil be stored for the sitemap
// Url is a string representing the url, Node is a pointer to the corresponding Node
// which represent that URL and skip is a boolean that is used to indicate if this URL should
// be followed (skip=false) or not

type Urlsstore struct {
	Url  string
	Node *Node
	Skip bool
}

// mutex to control Access to map
var Access sync.Mutex

// regexp to detect link with parameter and anchor I chose to not follow that link in the web crawler mechanism.
// It could be a configurable parameter
var q_regexp = regexp.MustCompile(`.*[\?#].*`)

// store Urls already visited to avoid cycles (reseted after each execution)
//var VisitedURLs = make(map[string]bool)

// this function aims to build a sitemap for  the url passed in argument (url.URL field)
// in collect all links of the considered page (make some checks) and return these links as a list of URLs
// to be treated as well. if links not strat with the based url pattern (e.g  the root url we provided) then it
// is considered as a leaf and we not follow this link (ask in the exercise)

func Sitemap(url *Urlsstore, baseurl *url.URL, alreadyScanned map[string]bool) ([]*Urlsstore, error) {
	var links []*Urlsstore
	// debug purpose only log, should be removed if we don't need it
	//fmt.Printf("DEBUG: Sitemap  %s, skip=%v\n", url.Url, url.Skip)

	Access.Lock() // take a lock to avoid concurrent  map read/write
	if !alreadyScanned[url.Url] {
		alreadyScanned[url.Url] = true
	} else {
		Access.Unlock()
		return nil, nil
	}
	Access.Unlock()

	// check if the url has the expected form http[s]://
	// that's a design choice. We could refine this check
	if !utils.IsCorrrectURL(url.Url) {
		return links, nil
	}

	// check if the url is under the initial base url
	if !utils.IsUnderBaseUrl(baseurl.String(), url.Url) {
		return links, nil
	}

	// skip scan if url is a leaf
	if url.Skip {
		return links, nil
	}

	// get the content of the page for the considered url
	resp, err := http.Get(url.Url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected return code when trying to get url: %s (%d)", url.Url, resp.StatusCode)
	}

	// We don't try follow links with mime type <> text/html
	// It's a choice and is could maybe be refined
	if !utils.IsTextURL(resp) {
		return nil, nil
	}

	//  get and return all links of that page
	list := getLinks(resp.Body, url, baseurl, alreadyScanned)

	return list, nil
}

func getLinks(body io.Reader, parent *Urlsstore, baseUrl *url.URL, alreadyScanned map[string]bool) []*Urlsstore {
	// debug purpose only log, should be removed if we don't need it
	//fmt.Printf("DEBUG: getLinks  %s\n", parent.Url)

	var links []*Urlsstore
	z := html.NewTokenizer(body)
	for {
		token := z.Next()

		switch token {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, a := range token.Attr {
					if a.Key == "href" { // we are interested only by http links
						u, err := url.Parse(a.Val)
						if err != nil || u == nil {
							log.Printf("Error parsing url %s (%v, %v)\n", a.Val, err, u)
							continue
						}
						//treatment of relative path
						lurl := baseUrl.ResolveReference(u) // we try to get the full url
						// clean URL path
						utils.CleanPath(lurl)
						lta := lurl.String()

						Access.Lock() // lock to avoid map concurrent access
						if !alreadyScanned[lta] {
							child := Node{Url: lta}
							parent.Node.Children = append(parent.Node.Children, &child)
							us := &Urlsstore{Url: lta, Node: &child, Skip: false}
							if q_regexp.MatchString(lta) {
								us.Skip = true
							} else {
								us.Skip = false
							}
							links = append(links, us)
						}
						Access.Unlock()
					}
				}
			}

		}
	}
	return links // normally this code will not be reached , but just in case
}

// helpers function to print the computed sitemap stored in root Node data struct

func PrintNode(node *Node, n int) {

	for i := 0; i < n; i++ {
		fmt.Print("  ")
	}
	if n > 0 {
		fmt.Print("- ")
	}
	fmt.Println(node.Url)
	for _, c := range node.Children {
		PrintNode(c, n+1)
	}
}
