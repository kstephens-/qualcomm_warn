package scraper

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
	"time"
)

// get html response
func getHttp(url string) (resp *http.Response, ok bool) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to retrieve \"" + url + "\"")
		ok = false
	} else {
		ok = true
	}
	return
}

type QualcommEvent struct {
	name   string
	date   time.Time
	detail string
}

func (q *QualcommEvent) String() string {
	return fmt.Sprintf("%#v", q)
}

func (q *QualcommEvent) SetAttributeInOrder(v string) {
	switch {
	case q.date.IsZero():
		q.date, _ = time.Parse("1/02/06", v)
	case q.name == "":
		q.name = v
	default:
		q.detail = v
	}
}

func ExtractEvents(http *http.Response) (events []QualcommEvent) {
	body := http.Body
	defer body.Close()

	events = extract(body)
	return
}

// extract qualcomm events
func extract(body io.Reader) (events []QualcommEvent) {

	text := html.NewTokenizer(body)
	event := QualcommEvent{}
	tr := false
	for {
		element := text.Next()
		if event.detail != "" {
			events = append(events, event)
			event = QualcommEvent{}
		}
		switch {
		case element == html.ErrorToken:
			// end of document
			return
		case element == html.StartTagToken:
			// start of table row
			token := text.Token()
			if token.Data == "tr" {
				tr = true
			}
		case element == html.EndTagToken:
			// end of table row
			token := text.Token()
			if token.Data == "tr" {
				tr = false
			}
		case element == html.TextToken && tr == true:
			// text nodes in table row
			d := (string)(text.Text())
			if len(d) > 1 {
				value := strings.TrimSpace(d)
				event.SetAttributeInOrder(value)
			}
		}
	}
	return
}

func main() {
	url := "https://www.sandiego.gov/qualcomm/event"
	resp, ok := getHttp(url)
	if !ok {
		fmt.Println("ERROR: Unable to process \"" + url + "\"")
	}
	events := ExtractEvents(resp)
	fmt.Println(events)
}
