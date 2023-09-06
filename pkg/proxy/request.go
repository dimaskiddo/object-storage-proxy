package proxy

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
)

func (osp *ObjectStorageProxy) objectStorageProxyBuilder(req *http.Request, accessKey string, region string) (*http.Request, error) {
	proxyURL := *req.URL
	endpointDomain := osp.Endpoint

	switch {
	case osp.UpstreamStyle == "virtual" && osp.LocalStyle == "path":
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

	case osp.UpstreamStyle == "path" && osp.LocalStyle == "virtual":
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

	proxyURL.Scheme = osp.Scheme
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

	if !osp.IsPublic {
		if err := osp.sign(osp.signer(accessKey, osp.SecretKey), proxyReq, region); err != nil {
			return nil, err
		}
	}

	copyHeaderData(req.Header, proxyReq.Header)
	return proxyReq, nil
}

func (osp *ObjectStorageProxy) objectStorageProxyRequest(req *http.Request) (*http.Request, error) {
	var accessKey, region string

	if !osp.IsPublic {
		accessKey, region, err := osp.validateHeaders(req)
		if err != nil {
			return nil, err
		}

		fakeReq, err := osp.fakeIncomingRequest(req, accessKey, region)
		if err != nil {
			return nil, err
		}

		compareAuthorization := subtle.ConstantTimeCompare([]byte(fakeReq.Header["Authorization"][0]), []byte(req.Header["Authorization"][0]))
		if compareAuthorization == 0 {
			fakeDumpReq, _ := httputil.DumpRequest(fakeReq, false)
			log.Println(log.LogLevelError, "Fake Dump Request: "+string(fakeDumpReq))

			intialDumpReq, _ := httputil.DumpRequest(req, false)
			log.Println(log.LogLevelError, "Initial Dump Request: "+string(intialDumpReq))

			return nil, fmt.Errorf("Invalid Signature in Authorization Header")
		}
	}

	if osp.Region != "" {
		region = osp.Region
	}

	if osp.Verbose {
		intialDumpReq, _ := httputil.DumpRequest(req, false)
		log.Println(log.LogLevelDebug, "Initial Dump Request: "+string(intialDumpReq))
	}

	proxyReq, err := osp.objectStorageProxyBuilder(req, accessKey, region)
	if err != nil {
		return nil, err
	}

	proxyReq.ContentLength = req.ContentLength

	if osp.Verbose {
		proxyDumpReq, _ := httputil.DumpRequest(proxyReq, false)
		log.Println(log.LogLevelDebug, "Proxy Dump Request: "+string(proxyDumpReq))
	}

	return proxyReq, nil
}
