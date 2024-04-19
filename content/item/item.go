// Package item provides the default functionality to Ponzu's entities/data types,
// how they interact with the API, and how to override or enhance their abilities
// using various interfaces.
package item

import (
	"net/http"
	"time"
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

// CSVFormattable is implemented with the method FormatCSV, which must return the ordered
// slice of JSON struct tag names for the type implementing it
type CSVFormattable interface {
	FormatCSV() []string
}

// Hookable provides our user with an easy way to intercept or add functionality
// to the different lifecycles/events a struct may encounter. Item implements
// Hookable with no-ops so our user can override only whichever ones necessary.
type Hookable interface {
	BeforeAPIResponse(http.ResponseWriter, *http.Request, interface{}) (interface{}, error)
	AfterAPIResponse(http.ResponseWriter, *http.Request, interface{}) error

	BeforeAPICreate(http.ResponseWriter, *http.Request) error
	AfterAPICreate(http.ResponseWriter, *http.Request) error

	BeforeAPIUpdate(http.ResponseWriter, *http.Request) error
	AfterAPIUpdate(http.ResponseWriter, *http.Request) error

	BeforeAPIDelete(http.ResponseWriter, *http.Request) error
	AfterAPIDelete(http.ResponseWriter, *http.Request) error

	BeforeAdminCreate(http.ResponseWriter, *http.Request) error
	AfterAdminCreate(http.ResponseWriter, *http.Request) error

	BeforeAdminUpdate(http.ResponseWriter, *http.Request) error
	AfterAdminUpdate(http.ResponseWriter, *http.Request) error

	BeforeAdminDelete(http.ResponseWriter, *http.Request) error
	AfterAdminDelete(http.ResponseWriter, *http.Request) error

	BeforeSave(http.ResponseWriter, *http.Request) error
	AfterSave(http.ResponseWriter, *http.Request) error

	BeforeDelete(http.ResponseWriter, *http.Request) error
	AfterDelete(http.ResponseWriter, *http.Request) error

	BeforeReject(http.ResponseWriter, *http.Request) error
	AfterReject(http.ResponseWriter, *http.Request) error

	BeforeEnable(http.ResponseWriter, *http.Request) error
	AfterEnable(http.ResponseWriter, *http.Request) error

	BeforeDisable(http.ResponseWriter, *http.Request) error
	AfterDisable(http.ResponseWriter, *http.Request) error
}

// Item should only be embedded into entities type structs.
type Item struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Timestamp int64  `json:"timestamp"`
	Updated   int64  `json:"updated"`
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

// BeforeAPIResponse is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAPIResponse(res http.ResponseWriter, req *http.Request, data interface{}) (interface{}, error) {
	return data, nil
}

// AfterAPIResponse is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAPIResponse(res http.ResponseWriter, req *http.Request, data interface{}) error {
	return nil
}

// BeforeAPICreate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAPICreate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterAPICreate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAPICreate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeAPIUpdate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAPIUpdate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterAPIUpdate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAPIUpdate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeAPIDelete is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAPIDelete(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterAPIDelete is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAPIDelete(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeAdminCreate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAdminCreate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterAdminCreate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAdminCreate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeAdminUpdate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAdminUpdate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterAdminUpdate is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAdminUpdate(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeAdminDelete is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeAdminDelete(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterAdminDelete is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterAdminDelete(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeSave is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeSave(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterSave is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterSave(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeDelete is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeDelete(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterDelete is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterDelete(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeReject is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeReject(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterReject is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterReject(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeEnable is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeEnable(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterEnable is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterEnable(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// BeforeDisable is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) BeforeDisable(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// AfterDisable is a no-op to ensure structs which embed Item implement Hookable
func (i *Item) AfterDisable(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// IndexContent determines if a type should be indexed for searching
// partially implements search.Searchable
func (i *Item) IndexContent() bool {
	return false
}
