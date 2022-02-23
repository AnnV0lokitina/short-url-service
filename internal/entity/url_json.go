package entity

import (
	"encoding/json"
)

type jsonURL struct {
	Checksum string `json:"checksum"`
	FullURL  string `json:"full_url"`
}

func (u *URL) MarshalJSON() ([]byte, error) {
	tmpURLInfo := jsonURL{
		Checksum: u.checksum,
		FullURL:  u.full,
	}

	return json.Marshal(tmpURLInfo)
}

func (u *URL) UnmarshalJSON(data []byte) error {
	var tmpURLInfo jsonURL
	if err := json.Unmarshal(data, &tmpURLInfo); err != nil {
		return err
	}

	u.checksum = tmpURLInfo.Checksum
	u.full = tmpURLInfo.FullURL

	return nil
}
