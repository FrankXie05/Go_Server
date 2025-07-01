package gologin

import "EMU_server/types"

type ProfileConfig struct {
	Name        string          `json:"name,omitempty"`
	OS          string          `json:"os,omitempty"`
	BrowserType string          `json:"browserType"`
	Proxy       types.Proxy     `json:"proxy"`
	Navigator   types.Navigator `json:"navigator"`
	Fonts       FontsConfig     `json:"fonts"`
}

type FontsConfig struct {
	EnableDomRect bool     `json:"enableDomRect"`
	EnableMasking bool     `json:"enableMasking"`
	Families      []string `json:"families"`
}

func NewProfileConfig(profileName string, proxy types.Proxy, navigator types.Navigator) *ProfileConfig {
	if proxy.Host == "" {
		panic(" proxy.Host 不能为空")
	}
	if navigator.UserAgent == "" {
		panic(" navigator.UserAgent 不能为空")
	}
	return &ProfileConfig{
		Name:        profileName,
		OS:          "win",
		BrowserType: "chrome",
		Proxy:       proxy,
		Navigator:   navigator,
		Fonts: FontsConfig{ // 设置默认字体
			EnableDomRect: true,
			EnableMasking: true,
			Families:      []string{"Arial", "Verdana", "Times New Roman"},
		},
	}
}
