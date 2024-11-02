package utils

import (
	"encoding/json"
	"os"
	"syscall"
)

func getSystemLanguage() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getUserDefaultUILanguage := kernel32.NewProc("GetUserDefaultUILanguage")

	langID, _, _ := getUserDefaultUILanguage.Call()
	primaryLangID := langID & 0x3ff

	switch primaryLangID {
	case 0x09:
		return "en"
	case 0x15:
		return "pl"
	default:
		return "en"
	}
}

func loadLangFile(lang string) (map[string]string, error) {
	var filename string = "assets/lang/en.json"

	if lang == "pl" {
		filename = "assets/lang/pl.json"
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var langData map[string]string
	err = json.Unmarshal(data, &langData)
	if err != nil {
		return nil, err
	}
	return langData, nil
}

func LoadLang() (map[string]string, error) {
	return loadLangFile(getSystemLanguage())
}
