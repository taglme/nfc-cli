package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	//"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/f2prateek/train"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

var midPath = "/licenses/mid"

type Signer struct {
	appID  string
	secret *rsa.PrivateKey
	cert   string
	mid    string
}

//train.Interceptor
func NewSigner(appID string, secret *rsa.PrivateKey, cert string) *Signer {
	return &Signer{
		appID:  appID,
		secret: secret,
		cert:   cert,
	}
}

func (s *Signer) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	if s.mid == "" {
		//mid not set, we should init it
		midReqPath := "http://" + req.URL.Host + midPath
		resp, err := http.Get(midReqPath)
		if err != nil {
			return nil, errors.Wrap(err, "Could not make request to get MID")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Could not read MID request body")
		}
		var midRes models.LicenseMID
		err = json.Unmarshal(body, &midRes)
		if err != nil {
			return nil, errors.Wrap(err, "Could not marshal MID request body")
		}
		s.mid = midRes.MID
	}
	err := s.Sign(req)
	if err != nil {
		return nil, errors.Wrap(err, "Could not sign request")
	}
	return chain.Proceed(req)
}

type DesktopClaims struct {
	Certs []string `json:"x5c"`
	jwt.StandardClaims
	refAud string
}

func (s *Signer) SetMID(mid string) {
	s.mid = mid
}

func (s *Signer) MID() string {
	return s.mid
}

func (s *Signer) Sign(req *http.Request) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	claims := DesktopClaims{
		[]string{s.cert},
		jwt.StandardClaims{
			Issuer:    s.appID,
			Audience:  "desktop:" + s.mid,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		"",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(s.secret)
	if err != nil {
		return fmt.Errorf("Request sign error: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+ss)
	return nil
}

func PrivateRSAKeyFromB64String(str string) (*rsa.PrivateKey, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("Key base64 decode error: %s", err.Error())
	}
	pKey, err := x509.ParsePKCS1PrivateKey(decodedKey)
	if err != nil {
		return nil, fmt.Errorf("Key parse error: %s", err.Error())
	}
	return pKey, nil
}
