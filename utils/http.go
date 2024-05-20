package utils

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/motemen/go-loghttp"
	"github.com/rs/zerolog/log"
)

var (
	LoggingHttpTransport = &loghttp.Transport{
		LogRequest: func(req *http.Request) {
			log.Trace().Msgf("[%p] %s %s %v", req, req.Method, req.URL, req.Header)
		},
		LogResponse: func(resp *http.Response) {
			log.Trace().Msgf("[%p] %d %s", resp.Request, resp.StatusCode, resp.Request.URL)
		},
	}
)

func URLJoin(base string, paths ...string) string {
	p := path.Join(paths...)
	out := fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
	// make sure we dont change the trailing slash
	if strings.HasSuffix(paths[len(paths)-1], "/") {
		out += "/"
	}
	return out
}

func MakeHTTPS(url string) string {
	url = strings.TrimPrefix(url, "https:")
	url = strings.TrimPrefix(url, "http:")
	url = strings.TrimPrefix(url, "//")

	return "https://" + url
}
