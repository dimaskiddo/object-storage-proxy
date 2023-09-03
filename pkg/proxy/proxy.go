package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
)

type Handler struct {
	Scheme        string
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Region        string
	UpstreamStyle string
	LocalStyle    string
	Verbose       bool
}

func NewObjectStorageProxy(scheme string, endpoint string, accessKey string, secretKey string, region string, upstreamStyle string, localStyle string, verbose bool) (*Handler, error) {
	handler := &Handler{
		Scheme:        scheme,
		Endpoint:      endpoint,
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Region:        region,
		UpstreamStyle: upstreamStyle,
		LocalStyle:    localStyle,
		Verbose:       verbose,
	}

	return handler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyReq, err := h.objectStorageProxyRequest(r)

	if err != nil {
		log.Println(log.LogLevelError, "Unable to Proxy Object Storage Request")
		w.WriteHeader(http.StatusBadRequest)

		if h.Verbose {
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
