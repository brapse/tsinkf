# tsinkf [![Build Status][5][6]]
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
[5]: https://secure.travis-ci.org/brapse/tsinkf.png
[6]: https://travis-ci.org/brapse/tsinkf

## Usage
Execute `tsinkf run` followed by a command.

Running:
```bash
$ tsinkf run wc -l /usr/share/dict/words
2013-04-24 22:43:58     NEW->RUNNING    wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
2013-04-24 22:43:58     RUNNING->SUCCEEDED      wc -l /usr/share/dict/words     d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz 
```

Performing the command again will result in no output because the input domain has been executed successfully.
Behind the scenes tsink will persist the result of the commands to disk (.tsinkf/ by default).
```{bash}
$ tsinkf run wc -l /usr/share/dict/words
```

`tsinkf show` can be used to inspect the state (.tsinkf/ by default).
The output contains the completion time, the state, the command and the
jobID (base64 version of the command).
```{bash}
$ tsinkf show
2013-04-24 22:43:58     SUCCEEDED       wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
```

A specific job can be inspected and output shown by including the `-v`
flag along with the jobID as parameters to `tsinkf show`.
```{bash}
$ tsinkf show -v d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
2013-04-24 22:43:58     SUCCEEDED     wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
235886 /usr/share/dict/words
```

Running `tsinkf reset` will reset the state of all the jobs. Subsequent executions will
append their output to the original output.
```{bash}
$ tsinkf reset -v
2013-04-24 22:48:25     SUCCEEDED->NEW  wc -l /usr/share/dict/words d2MgLWwgL3Vzci9zaGFyZS9kaWN0L3dvcmRz
```

## Status
* Beta quality
* Not used in production
* Some features and sketches and not fully fleshed out

## Todo
* tsink reset -hard  #=> delete the contents
* Redo the help and subcommand listing
* Refactor output redirecting

## License
BSD 2-Clause, see [LICENSE][4] for more details.

[4]: https://github.com/brapse/tsinkf/blob/master/LICENSE
