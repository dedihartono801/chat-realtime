package dto

import "time"

type CustomTime struct {
	time.Time
}

const layout = "2006-01-02T15:04:05.999999"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	if b == nil || string(b) == "null" {
		return nil
	}
	parsedTime, err := time.Parse(`"`+layout+`"`, string(b))
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type RecordInfo struct {
	CreatedAt   CustomTime `json:"created_at"`
	UpdatedAt   CustomTime `json:"updated_at"`
	DeletedAt   CustomTime `json:"deleted_at"`
	ID          int64      `json:"id"`
	UserChatID  int64      `json:"user_chat_id"`
	MessageText string     `json:"message_text"`
}

type InsertInfo struct {
	Columns         []ColumnInfo `json:"columns"`
	CommitTimestamp time.Time    `json:"commit_timestamp"`
	Errors          interface{}  `json:"errors"` // You may want to define a specific error struct here
	Record          RecordInfo   `json:"record"`
	Schema          string       `json:"schema"`
	Table           string       `json:"table"`
	Type            string       `json:"type"`
}
