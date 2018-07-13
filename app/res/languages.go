package res

type lang struct {
	Code string
	Name string
	Flag string
}

var Languages = []lang{{
	Code: "en-US",
	Name: "English",
	Flag: "us",
}, {
	Code: "es-LA",
	Name: "Español",
	Flag: "mx",
}, {
	Code: "zh-CN",
	Name: "中文(简体)",
	Flag: "cn",
}, {
	Code: "ja-JP",
	Name: "日本語",
	Flag: "jp",
}, {
	Code: "fr-FR",
	Name: "français",
	Flag: "fr",
}, {
	Code: "sv-SE",
	Name: "svenska",
	Flag: "se",
}, {
	Code: "ko-KR",
	Name: "한국어",
	Flag: "kr",
}, {
	Code: "el-GR",
	Name: "Ελληνικά",
	Flag: "gr",
}, {
	Code: "pl-PL",
	Name: "Polski",
	Flag: "pl",
}, {
	Code: "pt-BR",
	Name: "Português",
	Flag: "br",
}, {
	Code: "cs-CZ",
	Name: "Čeština",
	Flag: "cz",
}, {
	Code: "nl-NL",
	Name: "Nederlands",
	Flag: "nl",
}}

func IsValidLang(code string) bool {
	for _, lang := range Languages {
		if code == lang.Code {
			return true
		}
	}
	return false
}
