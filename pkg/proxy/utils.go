package proxy

import (
	"net/http"
)

func copyHeaderData(src http.Header, dst http.Header) {
	for srcHeaderKey, srcHeaderValue := range src {
		if _, isHeaderSet := dst[srcHeaderKey]; !isHeaderSet {
			for _, srcHeaderMultiValue := range srcHeaderValue {
				dst.Add(srcHeaderKey, srcHeaderMultiValue)
			}
		}
	}
}
