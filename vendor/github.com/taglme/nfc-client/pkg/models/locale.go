package models

type Locale int

const (
	LocaleEn Locale = iota + 1
	LocaleRu
)

func StringToLocale(s string) (Locale, bool) {
	switch s {
	case LocaleEn.String():
		return LocaleEn, true
	case LocaleRu.String():
		return LocaleRu, true
	}
	return 0, false
}

func (locale Locale) String() string {
	names := [...]string{
		"unknown",
		"en",
		"ru",
	}

	if locale < LocaleEn || locale > LocaleRu {
		return names[0]
	}
	return names[locale]
}
