package constants

type LanguageType string

const (
	LanguageC    LanguageType = "c"
	LanguageJava LanguageType = "java"
	LanguageGo   LanguageType = "go"
	LanguageCPP  LanguageType = "cpp"
)

var SupportLanguages = []LanguageType{
	LanguageGo,
	LanguageC,
	LanguageCPP,
}
