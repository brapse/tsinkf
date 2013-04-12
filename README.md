tsinkf
======
Perform a command on a set of arguments exactly once.

Usage
=====
```
$ tsinkf run --from="find /bin -type f|head" --to="wc -l"
```

```
$ tsinkf show
2012-01-01 2:00:00	wc -l /bin/bash	 SUCCEEDED
2012-01-01 2:00:00	wc -l /bin/cat   SUCCEEDED
2012-01-01 2:00:00	wc -l /bin/chmod SUCCEEDED
...
```

```
$ tsinkf show -v
2012-01-01 2:00:00	wc -l /bin/bash	 SUCCEEDED #=> 31
2012-01-01 2:00:00	wc -l /bin/cat	 SUCCEEDED #=> 426
...
```

```
$ tsinkf reset
$ tsinkf show
...
```
