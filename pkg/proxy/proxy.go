package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
)

type ObjectStorageProxy struct {
	Scheme        string
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Region        string
	UpstreamStyle string
	LocalStyle    string
	Verbose       bool
}

func NewObjectStorageProxy(osp ObjectStorageProxy) (*ObjectStorageProxy, error) {
	objectStorageProxy := &ObjectStorageProxy{
		Scheme:        osp.Scheme,
		Endpoint:      osp.Endpoint,
		AccessKey:     osp.AccessKey,
		SecretKey:     osp.SecretKey,
		Region:        osp.Region,
		UpstreamStyle: osp.UpstreamStyle,
		LocalStyle:    osp.LocalStyle,
		Verbose:       osp.Verbose,
	}

	return objectStorageProxy, nil
}

func (osp *ObjectStorageProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyReq, err := osp.objectStorageProxyRequest(r)

	if err != nil {
		log.Println(log.LogLevelError, "Unable to Proxy Object Storage Request")
		w.WriteHeader(http.StatusBadRequest)

		if osp.Verbose {
			w.Write([]byte(err.Error()))
		}

		return
	}
	defer proxyReq.Body.Close()

	proxyURL := url.URL{Scheme: proxyReq.URL.Scheme, Host: proxyReq.Host}

	proxyReverse := httputil.NewSingleHostReverseProxy(&proxyURL)
	proxyReverse.FlushInterval = 1

	proxyReverse.ServeHTTP(w, proxyReq)
}
