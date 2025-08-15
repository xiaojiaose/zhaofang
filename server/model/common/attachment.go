package common

import (
	"database/sql/driver"
	"encoding/json"
)

type AttachmentMap map[string]AttachmentArr

type Attachment struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (c AttachmentMap) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *AttachmentMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type AttachmentArr []Attachment

func (c AttachmentArr) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *AttachmentArr) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

const (
	IdentityRental  IdentityType = "rental"   // 租户
	IdentityLand    IdentityType = "land"     // 业主
	IdentityCoOwner IdentityType = "co_owner" // 共有产权人
	IdentityAgent   IdentityType = "agent"    // 代理人

)

type IdentityType string
type IdentityTypes struct {
	Rental  bool
	Land    bool
	CoOwner bool
	Agent   bool
}

func (c IdentityTypes) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *IdentityTypes) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
