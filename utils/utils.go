package utils

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
	syswindows "golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

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

func BgColor() application.RGBA {
	if IsWindows11OrGreater() {
		return application.NewRGBA(0, 0, 0, 0)
	} else {
		if IsDarkMode() {
			return application.NewRGBA(32, 32, 32, 255)
		} else {
			return application.NewRGBA(245, 241, 248, 255)
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

func GetAccentColor() (string, error) {
	accentColor, err := getRawAccentColor()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("#%06X", accentColor&0xFFFFFF), nil
}
