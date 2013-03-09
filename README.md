ALF: No problem!
================

Intro
-----
Alf is a framework for running scheduled tasks on remote hosts. A task can be a script, a program, or just
a shell command. Alf communicates with remote hosts via SSH, and thus, it doesn't require any client-side
installation other than the sshd service. Command scheduling is done with a cron-like line, and results
are stored in a database, as well as hosts and tasks details. Any manipulation of the database, such as
adding, modifying, reading a host, check, results, etc. is done over an HTTP api. So clients for
manipulating the information and querying the status can be developed separately and for example, we
can have a CLI interface, an HTTP interface or a mobile app.

An example use case of Alf could be a monitoring/alerting system, or an easy way to centrally manage 
cron tasks in a large network of servers, etc.


The main goals of the project are:

* Be able to detect easily failing tasks in remote hosts.
* Be easy to use, add hosts, tasks, contacts, notifications of failures, etc.
* Be easily extensible via language agnostic script/programs.
* Focus on the remote execution core, and provide an API for management. This way, there could be
command line clients, web interfaces, external automations, etc.
* Scale decently.
* Any existing nagios plugin must be valid for use in an Alf task.

Concepts
--------
* **Host:** A Host represents a remote ssh server where we are going to run Tasks on.
* **Task:** A Task wraps the command line, environment variables, etc to be run on a Host.
* **Schedule:** A Schedule represents a set of points in time where a Task should be run 
* **SSH Connection:** A SSH connection represents the network connection between Alf and the remote Host. there
    are 1 or more SSH Sessions over a connection.
* **SSH Session:** A SSH Session wraps a channel inside a SSH Connection.
* **Node:** A node represents an instance of the Alf core, it stores information about a particular instance
    of an Alf core, such as the Hosts it will run tasks for.
* **Result:** A Result stores the outcome of a Task after it has been run on a server, the Result contains
    information such as the exit code, stdout and stderr output, etc. We catalog results as even OK or PROBLEM.
    You expect a certain output from a Task, if you get something different, then it is a problem.
* **Action:** An action represents a link between a Result, and a Task. If a Result meets an Action requeriment,
    then another Task will be fired. This could be used for example, for sending out alarms or react to 
    problems. 
* **Client:** Any tool or application using the HTTP API.


Current implementation
----------------------
The current implementation is just a proof of concept of the idea. 

It is implemented in [Go](http://golang.org) where there is one [goroutine](http://golang.org/doc/effective_go.html#goroutines) for each (Host,Task) pair. 

There is no more than only one single SSH connection to each Host. The connection is shared between all 
the Tasks targeted to that Host by using SSH [channels](http://www.ietf.org/rfc/rfc4254.txt), which we call
'Sessions' in Alf.

There is an embedded webserver providing some entry points for data manipulation.

For storage, we are using MongoDB, but we keep the database code abstracted in the db package, and we use mostly
key/value storage without strong relationships at the database level, so it should be easy to use any other
database backend like Cassandra, Riak, etc. but right now the best available driver for Go is the mongodb one.


So far, we have:

* An HTTP interface that supports xml/json, with a few entry points
* The core functionality for running remote Tasks over SSH channels, sharing the connection, although right
    now every check runs on every hosts at a hardcoded interval.
* Connections are established with user/password.
* A config package, for dealing with the core preferences such as port to listen or db configuration details.
* Extensive logging, with a very flexible, dedicated configuration file.

Things/ideas we'd like to implement:

* Unit tests.
* A command line client.
* Implement PKI authentication
* Improve the API to allow manipulation of all the models.
* Implement the scheduling with a cron-like syntax.
* Migrate to Cassandra as far as we are happy with a Go driver ([gocql](http://github.com/tux21b/gocql/) 
    looks promising).
* Implement the Node, Action, HostTask, TaskSchedule functionality
* Improve robustness when a connection fails
* Implement a basic web interface
* Migrate the current logging configuration format to Yaml
* Make something useful with nagios' performance data output, such as pipe it to an Action that could send
    it to graphite, ganglia, etc.
* Keep a 'last N status' record for every (task,host) pair, so it can be quickly queried for producing
    a dashboard or similar
* Implement APIKeys or authentication for Clients using the API.
* Keep some internal statistics about checks/min, total number of hosts,checks, etc.  
* Implement max retention (expiration) for Results storage

Development environment notes
-----------------------------
The current development is being done with Go tip, and we should stick to Go 1.1 as soon as it is released.

