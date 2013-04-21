tsinkf
======
Perform a command on a set of arguments exactly once.

Install
=======
Install [Go 1][1], either [from source][2] or [with a prepackaged binary][3]. Then,
```bash
$ go get github.com/brapse/tsinkf
```

[1]: http://golang.org
[2]: http://golang.org/doc/install/source
[3]: http://golang.org/doc/install


Usage
=====

tsinkf run with -from argument and -to commands will execute every line
of the output of from as the last argument of the to command

```
$ tsinkf run -from="find /bin -type f|head" -to="wc -l"
```

Behind the scene tsink will persist result of each of the commands to
disk (.tsinkf/ by default) to ensure it does things exactly once.


tsinkf show will inspect the state (.tsinkf/ by default).
```bash
$ tsinkf show
2013-04-18 15:37:58     SUCCEEDED       wc -l /bin/cat d2MgLWwgL2Jpbi9jYXQ=
2013-04-18 15:37:58     SUCCEEDED       wc -l /bin/chmod d2MgLWwgL2Jpbi9jaG1vZA==
...
```

The output contains the completion time, the state, the command and the
command id (base64 version of the command)

Running tsinkf show in verbose mode (-v flag) will include the stdout of
the commands execution
```bash
$ tsinkf show -v
2013-04-18 15:37:58     SUCCEEDED       wc -l /bin/cat d2MgLWwgL2Jpbi9jYXQ=
57 /bin/cat

2013-04-18 15:37:58     SUCCEEDED       wc -l /bin/chmod d2MgLWwgL2Jpbi9jaG1vZA==
49 /bin/chmod
...
```

Running tsinkf will reset state of all the jobs, making it possible to
re-run everything
```bash
$ tsinkf reset
...
```

STATUS
======
* Alpha quality
* Not used in production
* Some features and sketches and not fully fleshed out

NOTES
=====
Jobs are identified by base64 encoding the full command. The current
persistance mechanism creates files named after the jobID. In cases in
which the encoded jobID is longer than 255 charecters, tsink will fail
to create a file and panic.

TODO
====
* tsink reset -hard  #=> delete the contents
* Redo the help and subcommand listing
* Refactor output redirecting
* Refactor job storage

LICENSE
=======
BSD 2-Clause, see [LICENSE][4] for more details.

[4]: https://github.com/brapse/tsinkf/blob/master/LICENSE
