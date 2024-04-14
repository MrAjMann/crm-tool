package model

import (
	"time"
)

type NoteCategory string

const (
	InteractionNote         NoteCategory = "Interaction"
	FeedbackNote            NoteCategory = "Feedback"
	InternalObservationNote NoteCategory = "Internal Observation"
	FollowUpNote            NoteCategory = "Follow-Up"
	OtherNote               NoteCategory = "Other"
)

type Note struct {
	NoteId     int
	CustomerId *int
	LeadId     *int
	Category   NoteCategory
	AuthorId   int    // this will be implemented after the main func of crm is done and user auth is added
	AuthorName string // this will be implemented after the main func of crm is done and user auth is added
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
