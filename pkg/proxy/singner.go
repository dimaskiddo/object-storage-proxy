package proxy

import (
	"bytes"
	"io"
	"net/http"
	"time"

	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

func (h *Handler) sign(signer *v4.Signer, req *http.Request, region string) error {
	return h.signWithTime(signer, req, region, time.Now())
}

func (h *Handler) signWithTime(signer *v4.Signer, req *http.Request, region string, signTime time.Time) error {
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
