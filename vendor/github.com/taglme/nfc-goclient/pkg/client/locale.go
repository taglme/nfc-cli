package client

import (
	"net/http"

	"github.com/f2prateek/train"
)

func NewLocalizer(locale string) train.Interceptor {
	return &Localizer{
		locale: locale,
	}
}

type Localizer struct {
	locale string
}

func (i *Localizer) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	req.Header.Set("Accept-Language", i.locale)
	return chain.Proceed(req)
}
