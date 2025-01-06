// Package item provides the default functionality to Ponzu's entities/data types,
// how they interact with the API, and how to override or enhance their abilities
// using various interfaces.
package item

import (
	"time"

	"github.com/fanky5g/ponzu/content/workflow"
)

// Readable enables an entity to have a Title property
type Readable interface {
	GetTitle() string
}

type Sluggable interface {
	Readable
	SetSlug(string)
	ItemSlug() string
}

// Identifiable enables a struct to have its ID set/get. Typically, this is done
// to set an ID to -1 indicating it is new for DB inserts, since by default
// a newly initialized struct would have an ID of 0, the int zero-value, and
// BoltDB's starting key per bucket is 0, thus overwriting the first record.
type Identifiable interface {
	ItemID() string
	SetItemID(string)
}

// Sortable ensures data is sortable by time
type Sortable interface {
	Time() int64
	Touch() int64
}

type Temporal interface {
	CreatedAt() int64
	SetCreatedAt(time.Time)
	UpdatedAt() int64
	SetUpdatedAt(time.Time)
}

// Item should only be embedded into entities type structs.
type Item struct {
	ID            string         `json:"id"`
	WorkflowState workflow.State `json:"workflow_state"`
	Slug          string         `json:"slug"`
	Timestamp     int64          `json:"timestamp"`
	Updated       int64          `json:"updated"`
}

func (i *Item) SetState(state workflow.State) {
	i.WorkflowState = state
}

func (i *Item) GetState() workflow.State {
	return i.WorkflowState
}

// Time partially implements the Sortable interface
func (i *Item) Time() int64 {
	return i.Timestamp
}

// Touch partially implements the Sortable interface
func (i *Item) Touch() int64 {
	return i.Updated
}

// SetSlug sets the item's slug for its URL
func (i *Item) SetSlug(slug string) {
	i.Slug = slug
}

func (i *Item) CreatedAt() int64 {
	return i.Timestamp
}

func (i *Item) UpdatedAt() int64 {
	return i.Updated
}

func (i *Item) SetCreatedAt(t time.Time) {
	i.Timestamp = t.UnixMilli()
}

func (i *Item) SetUpdatedAt(t time.Time) {
	i.Updated = t.UnixMilli()
}

// ItemSlug sets the item's slug for its URL
func (i *Item) ItemSlug() string {
	return i.Slug
}

// ItemID gets the Item's ID field
// partially implements the Identifiable interface
func (i *Item) ItemID() string {
	return i.ID
}

// SetItemID sets the Item's ID field
// partially implements the Identifiable interface
func (i *Item) SetItemID(id string) {
	i.ID = id
}

// IndexContent determines if a type should be indexed for searching
// partially implements search.Searchable
func (i *Item) IndexContent() bool {
	return false
}
