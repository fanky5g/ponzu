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

// Identifiable enables a struct to have its ID set/get. Typically, this is done
// to set an ID to -1 indicating it is new for DB inserts, since by default
// a newly initialized struct would have an ID of 0, the int zero-value, and
// BoltDB's starting key per bucket is 0, thus overwriting the first record.
type Identifiable interface {
	ItemID() string
	SetItemID(string)
}

// ItemID gets the Item's ID field
// partially implements the Identifiable interface
func (item *Item) ItemID() string {
	return item.ID
}

// SetItemID sets the Item's ID field
// partially implements the Identifiable interface
func (item *Item) SetItemID(id string) {
	item.ID = id
}

type Sluggable interface {
	Readable
	SetSlug(string)
	ItemSlug() string
}

// SetSlug sets the item's slug for its URL
func (item *Item) SetSlug(slug string) {
	item.Slug = slug
}

// ItemSlug sets the item's slug for its URL
func (item *Item) ItemSlug() string {
	return item.Slug
}

func (item *Item) SetState(state workflow.State) {
	item.WorkflowState = state
}

func (item *Item) GetState() workflow.State {
	return item.WorkflowState
}

// Time partially implements the Sortable interface
func (item *Item) Time() int64 {
	return item.Timestamp
}

// Touch partially implements the Sortable interface
func (item *Item) Touch() int64 {
	return item.Updated
}

func (item *Item) CreatedAt() int64 {
	return item.Timestamp
}

func (item *Item) UpdatedAt() int64 {
	return item.Updated
}

func (item *Item) SetCreatedAt(t time.Time) {
	item.Timestamp = t.UnixMilli()
}

func (item *Item) SetUpdatedAt(t time.Time) {
	item.Updated = t.UnixMilli()
}

// IndexContent determines if a type should be indexed for searching
// partially implements search.Searchable
func (item *Item) IndexContent() bool {
	return false
}
