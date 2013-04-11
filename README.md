Sink
====
Perform a command on a command and log the results

Usage
=====
sink --from="ls -al" --to="wc -l"

TODO
====
* Parse command line options
* Execute remote processes
* Create directory structure

# STATE
We'll need some way to maintain state. This can be as simple as
directory structure with log files.

Let item = A single link of the ouput of --from command

Jobs can have 4 states
* NEW       -> ready to run
* RUNNING   -> In the process of running
* FAILED    -> Returned 1
* SUCCESS   -> Returned 0

* DONE will contain all cmd have have finished

# FILESYSTEM

A set of files should be created for the todo list

* On creation, each task should create a file.
* When running, a symlink to that file should be created under running
* When finished, a symlink should be created under either success or
  failure
* The contents of the file should be stdout/stderr of running the cmd

# Problem
The problem with using "from" items as filenames is that it forces the
"items" to be file safe. What happeneds if that is not the case?

# Reruns
* Delete the symlink under failed

Can we FORCE it to be the case? is it worth it?

  ./sink/jobs/foo
  ./sink/jobs/bar
  ./sink/jobs/baz

  ./sink/jobs/done/
  ./sink/jobs/running
  ./sink/jobs/failure
  ./sink/jobs/success

# API 
store.Get(item) => Job object
store.Set(item, state) => TRUE|FALSE
