package storage

import (
	"github.com/mailhog/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// MongoDB represents MongoDB backed storage backend
type MultiTenantMongoDB struct {
	database    *mgo.Database
}

// CreateMongoDB creates a MongoDB backed storage backend
func CreateMultiTenantMongoDB(uri, db, coll string) *MultiTenantMongoDB {
	log.Printf("Connecting to MongoDB: %s\n", uri)
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %s", err)
		return nil
	}
	err = session.DB(db).C(coll).EnsureIndexKey("created")
	if err != nil {
		log.Printf("Failed creating index: %s", err)
		return nil
	}
	return &MultiTenantMongoDB{
		database:    session.DB(db),
	}
}

// Store stores a message in MongoDB and returns its storage ID
func (m *MultiTenantMongoDB) Store(msg *data.Message, tenant string) (string, error) {
	err := m.database.C(tenant).Insert(m)
	if err != nil {
		log.Printf("Error inserting message: %s", err)
		return "", err
	}
	return string(msg.ID), nil
}

// Count returns the number of stored messages
func (m *MultiTenantMongoDB) Count(tenant string) int {
	c, _ := m.database.C(tenant).Count()
	return c
}

// Search finds messages matching the query
func (m *MultiTenantMongoDB) Search(kind, query string, start, limit int, tenant string) (*data.Messages, int, error) {
	messages := &data.Messages{}
	var count = 0
	var field = "raw.data"
	switch kind {
	case "to":
		field = "raw.to"
	case "from":
		field = "raw.from"
	}
	err := m.database.C(tenant).Find(bson.M{field: bson.RegEx{Pattern: query, Options: "i"}}).Skip(start).Limit(limit).Sort("-created").Select(bson.M{
		"id":              1,
		"_id":             1,
		"from":            1,
		"to":              1,
		"content.headers": 1,
		"content.size":    1,
		"created":         1,
		"raw":             1,
	}).All(messages)
	if err != nil {
		log.Printf("Error loading messages: %s", err)
		return nil, 0, err
	}
	count, _ = m.database.C(tenant).Find(bson.M{field: bson.RegEx{Pattern: query, Options: "i"}}).Count()

	return messages, count, nil
}

// List returns a list of messages by index
func (m *MultiTenantMongoDB) List(start int, limit int, tenant string) (*data.Messages, error) {
	messages := &data.Messages{}
	err := m.database.C(tenant).Find(bson.M{}).Skip(start).Limit(limit).Sort("-created").Select(bson.M{
		"id":              1,
		"_id":             1,
		"from":            1,
		"to":              1,
		"content.headers": 1,
		"content.size":    1,
		"created":         1,
		"raw":             1,
	}).All(messages)
	if err != nil {
		log.Printf("Error loading messages: %s", err)
		return nil, err
	}
	return messages, nil
}

// DeleteOne deletes an individual message by storage ID
func (m *MultiTenantMongoDB) DeleteOne(id, tenant string) error {
	_, err := m.database.C(tenant).RemoveAll(bson.M{"id": id})
	return err
}

// DeleteAll deletes all messages stored in MongoDB
func (m *MultiTenantMongoDB) DeleteAll(tenant string) error {
	_, err := m.database.C(tenant).RemoveAll(bson.M{})
	return err
}

// Load loads an individual message by storage ID
func (m *MultiTenantMongoDB) Load(id, tenant string) (*data.Message, error) {
	result := &data.Message{}
	err := m.database.C(tenant).Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Printf("Error loading message: %s", err)
		return nil, err
	}
	return result, nil
}
