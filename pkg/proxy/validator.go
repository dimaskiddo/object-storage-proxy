package proxy

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"regexp"
)

var awsAuthCredential = regexp.MustCompile("Credential=([a-zA-Z0-9]+)/[0-9]+/([^/]+-?[0-9]+?)/s3/aws4_request")

func (osp *ObjectStorageProxy) validateHeaders(req *http.Request) (string, error) {
	headerAmzDate := req.Header["X-Amz-Date"]
	if len(headerAmzDate) != 1 {
		return "", fmt.Errorf("X-Amz-Date Header is Missing or Set Multiple Times: %v", req)
	}

	headerAuthorization := req.Header["Authorization"]
	if len(headerAuthorization) != 1 {
		return "", fmt.Errorf("Authorization Header is Missing or Set Multiple Times: %v", req)
	}

	match := awsAuthCredential.FindStringSubmatch(headerAuthorization[0])
	if len(match) != 3 {
		return "", fmt.Errorf("Invalid Authorization Header: Credential Not Found: %v", req)
	}

	key := match[1]
	region := match[2]

	if subtle.ConstantTimeCompare([]byte(key), []byte(osp.AccessKey)) == 1 {
		return region, nil
	}

	return "", fmt.Errorf("Invalid Access Key in Credential: %v", req)
}
