package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"bitbucket.org/mischief/libauth"
)

type results struct {
	Items []items `json: "@items"`
}

type items struct {
	Title string
	Link string
	Snippet string
	Image image `json: "@image"`
}

type image struct {
	ContextLink string
	ThumbnailLink string
}
	
var (
	nmax = flag.Int("m", 50, "Number of results per query")
	related = flag.String("r", "", "Search for sites related to [url]")
	isearch =   flag.Bool("i", false, "Image search")
	itype =   flag.String("it", "", "Image type [clipart|face|lineart|news|photo]")
	isize =   flag.String("is", "", "Image size [huge|icon|large|medium|small|xlarge|xxlarge]")
	icolor =  flag.String("ic", "", "Image color [black|blue|brown|gray|green|orange|ping|purple|red|teal|white|yellow]")
	iscale =  flag.String("id", "", "Image scale [color|gray|mono]")
	ftype =   flag.String("f", "", "Filetype [bmp|gif|png|jpg|svg|pdf]")
	exact =   flag.String("e", "", "Match string exactly")
	exclude = flag.String("x", "", "Phrase to exclude")
	site =    flag.String("u", "", "Limit search to URL")
	safe =    flag.String("s", "off", "Safe search [active|high|medium|off]")
	snippet = flag.Bool("sn", false, "Include short description in results")
	thumb =   flag.Bool("t", false, "Include thumbnails")
)

func keys() (string, string, error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}
	key, err := libauth.Getuserpasswd("proto=pass service=gcli user=%s", u.Username)
	if err != nil {
		return "", "", err
	}
	cx, err := libauth.Getuserpasswd("proto=pass service=gcse user=%s", u.Username)
	if err != nil {
		return "", "", err
	}
	return key.Password, cx.Password, nil
}

func search(url string, re *results) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", "gcli (gzip)")
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	gr, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return json.NewDecoder(gr).Decode(&re)
}

func buildurl(key, cx string, start int) string {
	var opts strings.Builder

	query := strings.Join(flag.Args(), "+")
	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&start=%d&maxResults=10&cx=%s&q=%s", key, start, cx, query)

	switch {
	case *isearch:
		opts.WriteString("&searchType=image")
	case *itype != "":
		if *isearch {
			opts.WriteString("&imageType=")
			opts.WriteString(*itype)
		}
/*	case *related != "":
		opts.WriteString(
	case *isize != "":
		if *isearch {
			search = search.ImgSize(*isize)
		}
	case *icolor != "":
		if *isearch {
			search = search.ImgDominantColor(*icolor)
		}
	case *iscale != "":
		if *isearch {
			search = search.ImgColorType(*iscale)
		}

	case *exact != "":
		search = search.ExactTerms(*exact)
	case *exclude != "":
		search = search.ExcludeTerms(*exclude)
	case *site != "":
		search = search.SiteSearch(*site)
*/
	case *ftype != "":
		opts.WriteString("&fileType=")
		opts.WriteString(*ftype)
	}
	return url+opts.String()
}

func handle(r results, lines chan string) {
	for _, item := range r.Items {
		var line strings.Builder
		line.WriteString(fmt.Sprintf("%s %s", item.Title, item.Link))
		if *snippet {
			snip := strings.Replace(item.Snippet, "\n", " ", -1)
			line.WriteString(snip)
			line.WriteString("\n")
		}
		if *isearch {
			line.WriteString(fmt.Sprintf("%s %s", 
				item.Image.ContextLink, 
				item.Image.ThumbnailLink,
			))
		}
		lines <- line.String()	
	}
}
	
func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}
	key, cx, err := keys()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: There's no determinism to the output
	// but running this in the main routine is very slow on many systems

	// The API only gives us 10 results at a time
	// We need to loop through to our max results
	// And fetch each selection of results
	var wg sync.WaitGroup
	lines := make(chan string)
	for i := 0; i <= *nmax; i+=10 {
		wg.Add(1)
		go func(){
			defer wg.Done()
			url := buildurl(key, cx, i)
			r := results{}
			search(url, &r)
			handle(r, lines)
		}()
	}
	go func(){
		for line := range lines {
			fmt.Println(line)
		}
	}()	
	wg.Wait()
}
