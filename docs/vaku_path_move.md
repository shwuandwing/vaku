## vaku path move

Move a vault path from one location to another

### Synopsis

Moves a path from one location to another. This is equivalent to 'vaku path copy' followed
by 'vaku path delete (not destroy)' on the target.

Example:
  vaku path move secret/foo secret/bar

```
vaku path move [source folder] [target path] [flags]
```

### Options

```
  -h, --help   help for move
```

### Options inherited from parent commands

```
  -o, --format string   The output format to use. One of: "json", "text" (default "json")
```

### SEE ALSO

* [vaku path](vaku_path.md)	 - Contains all vaku path functions, does nothing on its own

###### Auto generated by spf13/cobra on 29-Oct-2018