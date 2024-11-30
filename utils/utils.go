package utils

import (
	"github.com/wailsapp/wails/v2/pkg/options"
	syswindows "golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var hexChars = []byte("0123456789ABCDEF")

func IsWindows11OrGreater() bool {
	major, _, build := syswindows.RtlGetNtVersionNumbers()
	return major >= 10 && build >= 22000
}

func IsDarkMode() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return false
	}

	return value == 0
}

func BgColor() *options.RGBA {
	if IsWindows11OrGreater() {
		return options.NewRGBA(0, 0, 0, 0)
	} else {
		if IsDarkMode() {
			return options.NewRGBA(32, 32, 32, 255)
		} else {
			return options.NewRGBA(255, 255, 255, 255)
		}
	}
}

func getRawAccentColor() (uint32, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\DWM`, registry.QUERY_VALUE)
	if err != nil {
		return 0, err
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("ColorizationColor")
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}

func byteToHex(b byte) string {
	return string([]byte{hexChars[b>>4], hexChars[b&0xF]})
}

func GetAccentColor() (string, error) {
	accentColor, err := getRawAccentColor()
	if err != nil {
		return "", err
	}
	r := byteToHex(byte(accentColor >> 16))
	g := byteToHex(byte(accentColor >> 8))
	b := byteToHex(byte(accentColor))

	return "#" + r + g + b, nil
}
