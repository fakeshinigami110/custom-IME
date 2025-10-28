package commands

// Config ساختار مشترک برای همه کامندها
type Config struct {
	ProjectName string
	IMEName     string
	Label       string
	Icon        string
	LangCode    string
	Description string
	ConfigFile  string
	ProjectDir  string
}