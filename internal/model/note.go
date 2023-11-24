package model

import (
	"time"
)

// NoteCategory defines the type of note.
type NoteCategory string

const (
	InteractionNote         NoteCategory = "Interaction"          // Notes from interactions like calls, emails
	FeedbackNote            NoteCategory = "Feedback"             // Customer feedback or opinions
	InternalObservationNote NoteCategory = "Internal Observation" // Internal observations and insights
	FollowUpNote            NoteCategory = "Follow-Up"            // Follow-up actions or reminders
	OtherNote               NoteCategory = "Other"                // Any other type of note
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
