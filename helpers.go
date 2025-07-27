package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-telegram/bot/models"
	"golang.org/x/net/html"
)

// Gets the first link from the chain: message -> formatted message -> forwarded messages
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

// Finds and returns a link from forwarded messages or rich text
func getUrlFromEntityMsg(entityMsg []models.MessageEntity) string {
	for _, msg := range entityMsg {
		if msg.URL != "" {
			return msg.URL
		}
	}
	return ""
}

// Extracts a reference from a string
func getUrlFromMessage(messageText string) string {
	// Regular expression for finding URL:
	// 1. Starts with http:// or https://
	// 2. Domain: letters, numbers, periods, hyphens
	// 3. Path: any characters except spaces and punctuation
	// 4. Ignores periods/commas at the end
	re := regexp.MustCompile(`https?://(?:www\.)?[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+(?:/[^\s.,!?;:"'<>(){}]*)?`)

	// Looking for the first match
	match := re.FindString(messageText)
	if checkValidateUrl(match) {
		return match
	}
	return ""
}

func checkValidateUrl(urlText string) bool {
	u, err := url.Parse(urlText)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}

func extractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("не удалось извлечь домен")
	}

	// Remove 'www.' if there is one
	domain := strings.TrimPrefix(host, "www.")

	// We take only the second level domain for some cases
	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		domain = strings.Join(parts[len(parts)-2:], ".")
	}

	return domain, nil
}

// Gets the title inside the <h1> tag by parsing the html text
func getFirstH1Text(url string) string {
	doc, err := getHtmlData(url)
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

	if h1Text == "" {
		h1Text, err = extractDomain(url)
		if err != nil {
			// TODO
		}

	}

	return strings.TrimSpace(h1Text) // Remove spaces at the beginning and end
}

// Separate logic for obtaining an HTML document
func getHtmlData(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO - logs
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

// Getting a composite key for sync.Map
func getCompositeSyncMapKey(update *models.Update) CompositeSyncMapKey {
	return CompositeSyncMapKey{
		TelegramId: update.Message.From.ID,
		MsgId:      update.Message.ID,
	}
}

// The function forms a single string from the input data with a list of user links
func getListMsg(urls map[int32]Link) (outputText string) {
	for id, link := range urls {
		outputText += fmt.Sprintf("%d: <a href=\"%s\">%s</a>\n", id, link.URL, link.Title)
	}
	return outputText
}
