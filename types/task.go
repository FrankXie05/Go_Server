package types

import (
	"path/filepath"
)

type TaskRequest struct {
	Type         string      `json:"type,omitempty"`
	TrackingLink string      `json:"trackingLink,omitempty"`
	Script       string      `json:"script,omitempty"`
	SoiMetaData  SoiMetaData `json:"soiMetaData"`
	Proxy        Proxy       `json:"proxy"`
	Navigator    Navigator   `json:"navigator"`
}
type SoiMetaData struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password,omitempty"`
	ZipCode     string `json:"zip_code,omitempty"`
	Address     string `json:"address,omitempty"`
	Province    string `json:"province,omitempty"`
	City        string `json:"city,omitempty"`
}

type Proxy struct {
	Mode     string `json:"mode,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Navigator struct {
	UserAgent      string `json:"userAgent,omitempty"`
	Resolution     string `json:"resolution,omitempty"`
	Platform       string `json:"platform,omitempty"`
	MaxTouchPoints string `json:"maxTouchPoints,omitempty"`
	Language       string `json:"language,omitempty"`
}

type Response struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type ProfileEnv struct {
	TmpDir               string
	ProfileID            string
	ProfilePath          string
	ProfileDefaultPath   string
	ProfileZipPath       string
	ProfileZipUploadPath string
}

func NewProfileEnv(tmpDir, profileID string) *ProfileEnv {
	if profileID == "" {
		return nil
	}
	profileBase := filepath.Join(tmpDir, profileID)
	return &ProfileEnv{
		TmpDir:               tmpDir,
		ProfileID:            profileID,
		ProfilePath:          profileBase,
		ProfileDefaultPath:   filepath.Join(profileBase, "Default"),
		ProfileZipPath:       filepath.Join(tmpDir, profileID+".zip"),
		ProfileZipUploadPath: filepath.Join(tmpDir, profileID+"_upload.zip"),
	}
}
