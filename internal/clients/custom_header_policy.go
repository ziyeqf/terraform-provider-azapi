package clients

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type customHeaderPolicy struct {
	headers http.Header
}

func NewCustomHeaderPolicy(headers http.Header) *customHeaderPolicy {
	return &customHeaderPolicy{
		headers: headers,
	}
}

var _ policy.Policy = customHeaderPolicy{}

func (p customHeaderPolicy) Do(req *policy.Request) (*http.Response, error) {
	for k, vs := range p.headers {
		for _, v := range vs {
			req.Raw().Header.Set(k, v)
		}
	}
	return req.Next()
}
