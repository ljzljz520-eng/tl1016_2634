package model

type Record struct {
	ID, AssetTag, Name, Location, Status, Owner string
	Slots                                       []string
	Version                                     int
}
type AuditEvent struct {
	ID, RecordID, Action, Actor, Detail string
	Sequence                            int
}
type Workflow struct {
	ID, RecordID, Stage, State string
	Steps                      []string
}
type Attachment struct {
	ID, RecordID, Filename, Digest string
	Size                           int64
}

type Query struct {
	Text, Status, Location string
	Limit                  int
}
type Result struct {
	Records []Record
	Total   int
	Next    string
}
type Decision struct {
	Allowed bool
	Reason  string
}
