package proxy

import (
	"bytes"
	"io"
	"net/http"
	"time"

	creds "github.com/aws/aws-sdk-go/aws/credentials"
	signer_v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

func (osp *ObjectStorageProxy) sign(signer *signer_v4.Signer, req *http.Request, region string) error {
	return osp.signWithTime(signer, req, region, time.Now())
}

func (osp *ObjectStorageProxy) signWithTime(signer *signer_v4.Signer, req *http.Request, region string, signTime time.Time) error {
	body := bytes.NewReader([]byte{})

	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		body = bytes.NewReader(bodyBytes)
	}

	_, err := signer.Sign(req, body, "s3", region, signTime)
	return err
}

func (osp *ObjectStorageProxy) generateSigner(accessKey string, secretKey string) *signer_v4.Signer {
	return signer_v4.NewSigner(creds.NewStaticCredentialsFromCreds(creds.Value{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
	}))
}
