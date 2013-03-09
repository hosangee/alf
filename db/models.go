package db

import ()

// DbModel is the interface every method dealing with models must require.
type DbModel interface {
	SetId(string)
	GetId() string
}

// BaseModel implements the minimum required fields and methods for any model
type BaseModel struct {
	Id string
}

func (self *BaseModel) GetId() string {
	return self.Id
}

func (self *BaseModel) SetId(id string) {
	self.Id = id
}

// Host represents a hosts across the system
type Host struct {
	//Id       bson.ObjectId `bson:"_id,omitempty" xml:"-" json:"-"`
	BaseModel
	Hostname string
	Port     int
	Username string
	Tasks    []string
	Password string `json:"-"`
}

// TaskResult implements the output of a task run
type TaskResult struct {
	//Id     bson.ObjectId `bson:"_id,omitempty" xml:"-" json:"-"`
	BaseModel
	Status int
	Signal string
	Msg    string
	Stdout string
	Host   *Host
	Task   *Task
}

// Task implements the configuration of a particular task
type Task struct {
	//Id   bson.ObjectId `bson:"_id,omitempty" xml:"-" json:"-"`
	BaseModel
	Name     string
	Cmd      string
	Schedule string
}

// Node represents a node running tasks.
type Node struct {
	BaseModel
	Name string
}

// HostCheck represent the assignment of a Task to a Host
type HostTask struct {
	BaseModel
	HostId  string
	CheckId string
}

// NodeHost represent the assignment of a Host to a Node
type NodeHost struct {
	BaseModel
	NodeId string
	HostId string
}
