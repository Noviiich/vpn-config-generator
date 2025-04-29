package telegram

var countryFlags = map[string]string{
	"US": "🇺🇸",
	"RU": "🇷🇺",
	"DE": "🇩🇪",
	"FR": "🇫🇷",
	"JP": "🇯🇵",
	"NL": "🇳🇱",
}

func GetCountryFlag(countryCode string) string {
	return countryFlags[countryCode]
}
