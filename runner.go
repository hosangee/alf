package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	log "code.google.com/p/log4go"
	"github.com/msurdi/alf/db"
	"time"
)

// A runner represents the association between a particular task and a host.
// Every runner should run in its own goroutine
type Runner struct {
	Host  *db.Host
	Task  *db.Task
	alfDb *db.DB
}

func NewRunner(host *db.Host, task *db.Task, alfDb *db.DB) *Runner {
	return &Runner{
		Host:  host,
		Task:  task,
		alfDb: alfDb,
	}
}

// Run a task on a host. The task will be run in it's own
// ssh session on a shared connection.
// If the task succeeds to complete, its result is stored in the database, if it fails
// (due to a network problem, etc) then the registered error is logged
func (self *Runner) run() {
	// The task will run once every minute
	c := time.Tick(1 * time.Minute)
	for _ = range c {
		log.Info("Running task: " + self.Host.Hostname + "/" + self.Task.Id)
		var buffer bytes.Buffer

		// Get a session to the host
		connection := GetHostConnection(self.Host)
		session, err := connection.GetSession()

		if err != nil {
			log.Error("Can't get session for " + self.Host.Hostname + ": " + err.Error())
		} else {

			session.Stdout = &buffer

			// Run the task
			if err = session.Run(self.Task.Cmd); err != nil {
				log.Error("Error running task  " + self.Task.Name + "@" + self.Host.Hostname + ": " + err.Error())
			}
			// Build result
			// TODO: This could be much better encapsulated
			result := db.TaskResult{
				Status: 0,
				Msg:    "",
				Signal: "",
				Stdout: buffer.String(),
				Host:   self.Host,
				Task:   self.Task,
			}
			if err != nil {
				result.Status = err.(*ssh.ExitError).ExitStatus()
				result.Msg = err.(*ssh.ExitError).Msg()
				result.Signal = err.(*ssh.ExitError).Signal()
			}
			err = self.alfDb.Results.Save(&result)
			if err != nil {
				log.Error("Error saving task result " + self.Task.Name + "@" + self.Host.Hostname + ": " + err.Error())
			}
		}
	}
}
