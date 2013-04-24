# tsinkf
Stateful command execution

## Motivation
Sometimes you want to perform command exactly
once. Tsinkf keeps track what commands have been executed and their
results

## Install
Install [Go 1][1], either [from source][2] or [with a prepackaged binary][3]. Then,
```bash
$ go get github.com/brapse/tsinkf
```

[1]: http://golang.org
[2]: http://golang.org/doc/install/source
[3]: http://golang.org/doc/install

## Usage
Execute `tsinkf run` followed by a command.

Running:
```bash
$ tsinkf run wc -l /usr/share/dict/words
2013-04-24 22:43:58     NEW->RUNNING    wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
2013-04-24 22:43:58     RUNNING->SUCCEEDED      wc -l /usr/share/dict/words     d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz 
```

followed by
```{bash}
$ tsinkf run wc -l /usr/share/dict/words
```

Will result in no output because the input domain has been executed successfully.
Behind the scene tsink will persist result of each of the commands to
disk (.tsinkf/ by default) to ensure it does things exactly once.

`tsinkf show` will inspect the state (.tsinkf/ by default).
```{bash}
$ tsinkf show
2013-04-24 22:43:58     SUCCEEDED       wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
```
The output contains the completion time, the state, the command and the
jobID (base64 version of the command).

A specific job can be inspected and output shown by including the `-v`
flag along with the jobID as parameters to `tsinkf show`.
```{bash}
2013-04-24 22:43:58     SUCCEEDED     wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
235886 /usr/share/dict/words
```


Running `tsinkf reset` state of all the jobs, making it possible to
re-run everything.
```{bash}
$ tsinkf reset -v
2013-04-24 22:48:25     SUCCEEDED->NEW  wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
```

## Status
* Beta quality
* Not used in production
* Some features and sketches and not fully fleshed out

## Notes
Jobs are identified by base64 encoding the full command. The current
persistance mechanism creates files named this jobID. In cases in
which the encoded jobID is longer than 255 charecters, tsink will fail
to create a file and crash.

## Todo
* tsink reset -hard  #=> delete the contents
* Redo the help and subcommand listing
* Refactor output redirecting

## License
BSD 2-Clause, see [LICENSE][4] for more details.

[4]: https://github.com/brapse/tsinkf/blob/master/LICENSE
