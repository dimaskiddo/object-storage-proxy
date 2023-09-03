package proxy

import (
	"net/http"
	"net/http/httputil"
	"strings"

	signer_v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
)

func (h *Handler) objectStorageProxyBuilder(signer *signer_v4.Signer, req *http.Request, region string) (*http.Request, error) {
	proxyURL := *req.URL
	endpointDomain := h.Endpoint

	switch {
	case h.UpstreamStyle == "virtual" && h.LocalStyle == "path":
		bucketName := ""

		urlSplit := strings.Split(req.URL.EscapedPath(), "/")
		if len(urlSplit) > 1 && urlSplit[1] != "" {
			bucketName = urlSplit[1]

			if len(urlSplit) > 2 {
				urlSplit = append(urlSplit[:1], urlSplit[2:]...)
			} else {
				urlSplit = urlSplit[:1]
			}
		}

		if bucketName != "" {
			endpointDomain = bucketName + "." + endpointDomain
		}

		proxyURL.Path = strings.Join(urlSplit, "/")

	case h.UpstreamStyle == "path" && h.LocalStyle == "virtual":
		bucketName := ""

		domainSplit := strings.Split(req.Host, ".")
		if len(domainSplit) > 1 {
			bucketName = domainSplit[0]
			domainSplit = domainSplit[1:]
		}

		endpointDomain = strings.Join(domainSplit, ".")
		if len(proxyURL.Path) > 1 {
			proxyURL.Path = "/" + bucketName + proxyURL.Path
		} else {
			proxyURL.Path = "/" + bucketName
		}
	}

	proxyURL.Scheme = h.Scheme
	proxyURL.Host = endpointDomain
	proxyURL.RawPath = req.URL.Path

	proxyReq, err := http.NewRequest(req.Method, proxyURL.String(), req.Body)
	if err != nil {
		return nil, err
	}

	if headerValue, isHeaderSet := req.Header["Content-Type"]; isHeaderSet {
		proxyReq.Header["Content-Type"] = headerValue
	}

	if headerValue, isHeaderSet := req.Header["Content-Md5"]; isHeaderSet {
		proxyReq.Header["Content-Md5"] = headerValue
	}

	if err := h.sign(signer, proxyReq, region); err != nil {
		return nil, err
	}

	copyHeaderData(req.Header, proxyReq.Header)
	return proxyReq, nil
}

func (h *Handler) objectStorageProxyRequest(req *http.Request) (*http.Request, error) {
	signer := h.Signer

	if h.Verbose {
		intialDumpReq, _ := httputil.DumpRequest(req, false)
		log.Println(log.LogLevelDebug, "Initial Dump Request: "+string(intialDumpReq))
	}

	proxyReq, err := h.objectStorageProxyBuilder(signer, req, h.Region)
	if err != nil {
		return nil, err
	}

	proxyReq.ContentLength = req.ContentLength

	if h.Verbose {
		proxyDumpReq, _ := httputil.DumpRequest(proxyReq, false)
		log.Println(log.LogLevelDebug, "Proxy Dump Request: "+string(proxyDumpReq))
	}

	return proxyReq, nil
}
