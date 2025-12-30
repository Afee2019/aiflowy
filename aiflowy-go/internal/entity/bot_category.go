package entity

import "time"

// BotCategory represents a bot category entity
type BotCategory struct {
	ID           int64     `json:"id,string"`
	CategoryName string    `json:"categoryName"`
	SortNo       int       `json:"sortNo"`
	Status       int       `json:"status"`
	Created      time.Time `json:"created"`
	CreatedBy    int64     `json:"createdBy,string"`
	Modified     time.Time `json:"modified"`
	ModifiedBy   int64     `json:"modifiedBy,string"`
}
