package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-telegram/bot/models"
	"golang.org/x/net/html"
)

// getFirstUrl Gets the first link from the chain: message -> formatted message -> forwarded messages
func getFirstUrl(urlMsgText string, urlEntitiesText ...[]models.MessageEntity) string {
	if urlText := getUrlFromMessage(urlMsgText); urlText != "" {
		return urlText
	}
	for _, urlEntText := range urlEntitiesText {
		if urlText := getUrlFromEntityMsg(urlEntText); urlText != "" {
			return urlText
		}
	}
	return ""
}

// getUrlFromMessage Extracts a reference from a string
func getUrlFromMessage(messageText string) string {
	match := regexpUrl(messageText, false)
	if match != "" && isUrl(match) {
		return match
	}
	return ""
}

// getUrlFromEntityMsg Finds and returns a link from forwarded messages or rich text
func getUrlFromEntityMsg(entityMsg []models.MessageEntity) string {
	for _, msg := range entityMsg {
		if msg.URL != "" {
			return msg.URL
		}
	}
	return ""
}

// regexpUrl Checks if url is in a string
func regexpUrl(messageText string, fullString bool) string {
	var re *regexp.Regexp
	if fullString {
		// checks that a string exactly matches a regular expression
		re = regexp.MustCompile(`^https://(?:www\.)?[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+/?([^:\s]*)?$`)
	} else {
		re = regexp.MustCompile(`https://(?:www\.)?[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+/?([^:\s]*)?`)
	}
	return re.FindString(messageText)
}

func isUrl(urlText string) bool {
	if regexpUrl(urlText, true) == "" {
		return false
	}
	u, err := url.Parse(urlText)
	return err == nil && u.Scheme != "" && u.Host != "" && u.Scheme == "https" && u.Port() == ""
}

func extractDomain(msgUrl string) string {
	u, err := url.Parse(msgUrl)
	if err != nil {
		log.Printf("Incorrect URL: %v, error: %v", msgUrl, err)
		return ""
	}

	host := u.Host

	// Remove "www." at the beginning
	domain := strings.TrimPrefix(host, "www.")

	return domain
}

func getTitle(msgUrl string) string {
	msgTitle := getFirstH1Text(msgUrl)
	if msgTitle == "" {
		msgTitle = extractDomain(msgUrl)
		if msgTitle == "" {
			msgTitle = msgUrl
		}
	}

	msgTitle = strings.TrimSpace(msgTitle) // Remove spaces at the beginning and end

	// Cut the string exactly by characters, without breaking them in the encoding
	maxTitleLen := 100
	if len([]rune(msgTitle)) > maxTitleLen {
		msgTitle = strings.TrimSpace(string([]rune(msgTitle)[:maxTitleLen]))
	}

	return msgTitle
}

// getFirstH1Text Gets the title inside the <h1> tag by parsing the html text
func getFirstH1Text(msgUrl string) string {
	doc, err := getHtmlData(msgUrl)
	if err != nil {
		return "" // If you didn't find the title - no problem, we'll just leave a link
	}

	// Find first <h1> element and extract all text inside it
	var h1Text string
	var findH1 func(*html.Node)
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			var extractText func(*html.Node)
			extractText = func(node *html.Node) {
				if node.Type == html.TextNode {
					h1Text += strings.TrimSpace(node.Data) + " "
				}
				// We collect all the text that may be in nested <h1> tags
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					extractText(c)
				}
			}
			extractText(n)
			return
		}
		// Recursively, inside each tag, we start searching for nested tags (up to the first <h1> tag)
		for c := n.FirstChild; c != nil && h1Text == ""; c = c.NextSibling {
			findH1(c)
		}
	}
	findH1(doc)

	return h1Text
}

// getHtmlData Separate logic for obtaining an HTML document
func getHtmlData(url string) (*html.Node, error) {
	// Take precautions when receiving html from a site
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return nil, fmt.Errorf("URL does not return HTML content")
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	return doc, nil
}

// getCompositeSyncMapKey Getting a composite key for sync.Map
func getCompositeSyncMapKey(update *models.Update) CompositeSyncMapKey {
	return CompositeSyncMapKey{
		TelegramId: update.Message.From.ID,
		MsgId:      update.Message.ID,
	}
}

// getListMsg The function forms a single string from the input data with a list of user links
func getListMsg(urls map[int32]Link) (outputText string) {
	keys := getSortKeys(urls)
	for _, id := range keys {
		linkData := urls[int32(id)]
		outputText += fmt.Sprintf("%d: <a href=\"%s\">%s</a>\n", id, linkData.URL, linkData.Title)
	}
	return outputText
}

// getSortKeys Sort keys in map in ascending order
func getSortKeys(unsortMap map[int32]Link) []int {
	keys := make([]int, 0, len(unsortMap))
	for k := range unsortMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys) // sort ascending
	return keys
}
