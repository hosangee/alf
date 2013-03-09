/*
	This module abstracts all the database related functionality. We expose the models and
	a single DB type which contains DAOs for every model we handle.
*/
package db

import (
	log "code.google.com/p/log4go"
	"github.com/tux21b/gocql/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// DB implements the database connection, and holds DAOs for every model we have.
type DB struct {
	// session is the mongodb session object
	session *mgo.Session

	//
	Hosts   *DAO
	Tasks   *DAO
	Results *DAO
}

func NewDB() *DB {
	return &DB{}
}

func NewId() (id string) {
	id = uuid.TimeUUID().String()
	return
}

// Connect connects to the database
func (self *DB) Connect(url string) (err error) {
	log.Info("Connecting to the database: %s", url)
	self.session, err = mgo.Dial(url)
	if err == nil {
		log.Info("Successfully connected")
		self.session.SetMode(mgo.Monotonic, true)
		self.Hosts = self.newDAO("hosts")
		self.Tasks = self.newDAO("tasks")
		self.Results = self.newDAO("results")
	} else {
		log.Error("Unable to connect to the database: %s", err.Error())
	}
	return
}

// Close closes the database connection
func (self *DB) Close() {
	log.Info("Closing db connection")
	self.session.Close()
}

// DAO abstracts the access to a model collection in the database
type DAO struct {
	namespace  string
	collection *mgo.Collection
}

// Returns a new DAO associated to collection
func (self *DB) newDAO(namespace string) (dao *DAO) {
	dao = &DAO{
		collection: self.session.DB("").C(namespace),
		namespace:  namespace,
	}
	return
}

// FindAll returns the complete list of models from the collection
// use with care, as it can take up all the available memory if there
// are no limits imposed by the db itself.
func (self *DAO) FindAll(entries interface{}) (err error) {
	log.Debug("FindAll: %s", self.namespace)
	err = self.collection.Find(bson.M{}).All(entries)
	if err != nil {
		log.Error("Error retrieving tasks: %s", err.Error())
	}
	return
}

// FindOneByField will store in '*entry' the first entry whose field 'field' matches
// the value 'value'
func (self *DAO) FindOneByField(field string, value interface{}, entry DbModel) (err error) {
	log.Debug("FindOneByField: field:%s,value:%s: %s", field, value, err.Error())
	err = self.collection.Find(bson.M{field: value}).One(entry)
	if err != nil {
		log.Error("Error:FindOneByField field:%s,value:%s: %s", field, value, err.Error())
	}
	return
}

// FindById is a shortcut for FindOneByField where field is "_id" and
func (self *DAO) FindById(id string, entry DbModel) (err error) {
	return self.FindOneByField("_id", bson.ObjectIdHex(id), entry)
}

// Save will save the provided 'entry' in the collection as a new entry.
// If you want to update an existing entry, use the Update() instead
func (self *DAO) Save(entry DbModel) (err error) {
	log.Debug("Save: %v", entry)
	entry.SetId(NewId())
	err = self.collection.Insert(entry)
	if err != nil {
		log.Error("Error:Save: %s", err.Error())
	}
	return
}
