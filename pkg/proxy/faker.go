package proxy

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	signer_v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

var awsAuthSignedHeader = regexp.MustCompile("SignedHeaders=([a-zA-Z0-9;-]+)")

func (osp *ObjectStorageProxy) fakeIncomingRequest(signer *signer_v4.Signer, req *http.Request, region string) (*http.Request, error) {
	fakeReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	fakeReq.URL.RawPath = req.URL.Path

	headerAuthorization := req.Header.Get("authorization")
	match := awsAuthSignedHeader.FindStringSubmatch(headerAuthorization)

	if len(match) == 2 {
		fakeReq.Header = req.Header
	}

	signTime, err := time.Parse("20060102T150405Z", req.Header["X-Amz-Date"][0])
	if err != nil {
		return nil, fmt.Errorf("error parsing X-Amz-Date %v - %v", req.Header["X-Amz-Date"][0], err)
	}

	if err := osp.signWithTime(signer, fakeReq, region, signTime); err != nil {
		return nil, err
	}

	return fakeReq, nil
}
