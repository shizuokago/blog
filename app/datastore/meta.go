package datastore

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Meta struct {
	Key       *datastore.Key `datastore:"-" json:"-"`
	Version   int            `datastore:",noindex" json:"version"`
	Deleted   bool           `json:"deleted"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type HasVersion interface {
	IncrementVersion()
	HasTime
}

type HasTime interface {
	SetTime()
	HasKey
}

type HasKey interface {
	GetKey() *datastore.Key
	SetKey(*datastore.Key)
}

func (m *Meta) GetKey() *datastore.Key {
	return m.Key
}

func (m *Meta) SetKey(key *datastore.Key) {
	m.Key = key
}

func (m *Meta) SetTime() {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()
}

func (m *Meta) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Meta) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *Meta) GetVersion() int {
	return m.Version
}

func (m *Meta) IncrementVersion() {
	m.Version++
}
