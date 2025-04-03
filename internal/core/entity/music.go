package entity

type Music struct {
	ID            string
	TimetableID   string
	Order         int
	Artist        string
	Title         string
	StreamAllowed bool
	Note          string
}
