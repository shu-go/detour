# detour

Windows shortcut replacer tool.



## Usage

### One line

```sh
$ detour -v -r oldserver:newserver -r oldpath:newpath  myserver.lnk
```

> detour  [-v]  -r OLD:NEW  [-r OLD:NEW]...  [--dry-run]  TARGET_PATTERNS...



### With rule set

```shell
$ cat myrules.txt
# server name
oldserver:newserver
# paths
oldpath:newpath

$ detour -v --rule-set myrules.txt  myserver.lnk
```

> detour  [-v]  --rule-set FILE_NAME  [--dry-run]  TARGET_PATTERNS...



### about TARGET_PATTERNS

TARGET_PATTERNS is glob patterns. You can use `**` to dive into directories.

Each pattern gets completed with `.lnk`.

