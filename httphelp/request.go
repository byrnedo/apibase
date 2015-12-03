package httphelp
import (
	"regexp"
	"net/http"
)

var (
	acceptsHtmlRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXmlRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJsonRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
)

func AcceptsJson(r *http.Request) bool {
	return acceptsJsonRegex.MatchString(r.Header("Accept"))
}

func AcceptsHtml(r *http.Request) bool {
	return acceptsHtmlRegex.MatchString(r.Header("Accept"))
}

func AcceptsXml(r *http.Request) bool {
	return acceptsXmlRegex.MatchString(r.Header("Accept"))
}
