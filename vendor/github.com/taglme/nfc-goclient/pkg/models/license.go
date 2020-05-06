package models

import (
	"time"

	"github.com/pkg/errors"
)

type License struct {
	ID           string
	Owner        string
	Email        string
	Machine      string
	Type         string
	Start        time.Time
	End          time.Time
	Support      time.Time
	Features     []string
	Applications []AppLicense
}

type LicenseResource struct {
	ID           string               `json:"id"`
	Owner        string               `json:"owner"`
	Email        string               `json:"email"`
	Machine      string               `json:"machine"`
	Type         string               `json:"type"`
	Start        string               `json:"start"`
	End          string               `json:"end"`
	Support      string               `json:"support"`
	Features     []string             `json:"features"`
	Applications []AppLicenseResource `json:"applications"`
}

type AppLicense struct {
	ID      string
	Name    string
	Type    string
	Start   time.Time
	End     time.Time
	Support time.Time
}
type AppLicenseResource struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Support string `json:"support"`
}
type LicenseMID struct {
	MID string `json:"mid"`
}

func (l AppLicense) IsActive() (ok bool) {
	if l.ID != "" {
		if l.End.IsZero() {
			ok = true
		} else {
			if time.Now().Before(l.End) {
				ok = true
			}
		}

	}
	return
}

func (l License) IsActive() (ok bool) {
	if l.ID != "" {
		if l.End.IsZero() {
			ok = true
		} else {
			if time.Now().Before(l.End) {
				ok = true
			}
		}

	}
	return
}

func (r *AppLicenseResource) ToAppLicense() (appLicense AppLicense, err error) {
	var licenseEnd, licenseStart, licenseSupport time.Time
	if r.End != "" {
		licenseEnd, err = time.Parse("2006-01-02", r.End)
		if err != nil {
			return appLicense, errors.Wrap(err, "Can't parse app license end time")
		}
	}
	if r.Start != "" {
		licenseStart, err = time.Parse("2006-01-02", r.Start)
		if err != nil {
			return appLicense, errors.Wrap(err, "Can't parse app license start time")
		}
	}
	if r.Support != "" {
		licenseSupport, err = time.Parse("2006-01-02", r.Support)
		if err != nil {
			return appLicense, errors.Wrap(err, "Can't parse app license support time")
		}
	}
	appLicense = AppLicense{
		ID:      r.ID,
		Name:    r.Name,
		Type:    r.Type,
		Start:   licenseStart,
		End:     licenseEnd,
		Support: licenseSupport,
	}
	return

}
func (r *LicenseResource) ToLicense() (license License, err error) {
	var licenseEnd, licenseStart, licenseSupport time.Time
	if r.End != "" {
		licenseEnd, err = time.Parse("2006-01-02", r.End)
		if err != nil {
			return license, errors.Wrap(err, "Can't parse license end time")
		}
	}
	if r.Start != "" {
		licenseStart, err = time.Parse("2006-01-02", r.Start)
		if err != nil {
			return license, errors.Wrap(err, "Can't parse license start time")
		}
	}
	if r.Support != "" {
		licenseSupport, err = time.Parse("2006-01-02", r.Support)
		if err != nil {
			return license, errors.Wrap(err, "Can't parse license support time")
		}
	}
	appLicenses := []AppLicense{}
	for _, appLicenseRes := range r.Applications {
		appLicense, err := appLicenseRes.ToAppLicense()
		if err != nil {
			return license, errors.Wrap(err, "Can't convert app license resource")
		}
		appLicenses = append(appLicenses, appLicense)
	}

	license = License{
		ID:           r.ID,
		Owner:        r.Owner,
		Email:        r.Email,
		Machine:      r.Machine,
		Type:         r.Type,
		Start:        licenseStart,
		End:          licenseEnd,
		Support:      licenseSupport,
		Features:     r.Features,
		Applications: appLicenses,
	}
	return

}
