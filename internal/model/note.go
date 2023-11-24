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
	NoteId     string
	LeadId     string
	Category   NoteCategory
	AuthorId   string
	AuthorName string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
