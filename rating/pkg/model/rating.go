package model

// RecordID defines a record id. Together with RecordType identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID identifies unique records across all types.
type RecordType string

// Existing record types.
const (
	RecordTypeMovie = RecordType("movie")
)

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for some record.
type Rating struct {
	RecordID   RecordID    `json:"record_id"`
	RecordType RecordType  `json:"record_type"`
	UserID     UserID      `json:"user_id"`
	Value      RatingValue `json:"value"`
}

// RatingEvent defines an event containing rating information.
type RatingEvent struct {
	UserID     UserID          `json:"user_id"`
	RecordID   RecordID        `json:"record_id"`
	RecordType RecordType      `json:"record_type"`
	Value      RatingValue     `json:"value"`
	ProviderID string          `json:"provider_id"`
	EventType  RatingEventType `json:"event_type"`
}

// RatingEventType defines a type of a rating event.
type RatingEventType string

const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
