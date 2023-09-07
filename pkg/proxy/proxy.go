package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	signer_v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

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
	IsPublic      bool
	Insecure      bool
	Verbose       bool
	Signer        *signer_v4.Signer
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
		IsPublic:      osp.IsPublic,
		Insecure:      osp.Insecure,
		Verbose:       osp.Verbose,
		Signer:        osp.NewSigner(osp.AccessKey, osp.SecretKey),
	}

	return objectStorageProxy, nil
}

func (osp *ObjectStorageProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyReq, err := osp.objectStorageProxyRequest(r)

	if err != nil {
		log.Println(log.LogLevelError, err.Error())

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
