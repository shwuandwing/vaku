## vaku

Vaku CLI extends the official Vault CLI with useful high-level functions

### Synopsis

Vaku CLI extends the official Vault CLI with useful high-level functions

The Vaku CLI is intended to be used side by side with the official Vault CLI,
and only provides functions to extend the existing functionality. Many of the 'vaku path'
functions are very similar (or even less featured) than the vault CLI equivalent. However
vaku works on both v1 and v2 secret mounts, and can even copy/move secrets between them.

Vaku does not log you in to vault or help you with getting a token. Like the CLI,
it will look for a token first at the VAULT_TOKEN env var and then in ~/.vault-token

Built by Sean Lingren <sean@lingrino.com>
CLI documentation is available using 'vaku help [cmd]'
API documentation is available at https://godoc.org/github.com/lingrino/vaku/vaku

### Options

```
  -o, --format string   The output format to use. One of: "json", "text" (default "json")
  -h, --help            help for vaku
```

### SEE ALSO

* [vaku folder](vaku_folder.md)	 - Contains all vaku folder functions, does nothing on its own
* [vaku path](vaku_path.md)	 - Contains all vaku path functions, does nothing on its own
* [vaku version](vaku_version.md)	 - Returns the current Vaku CLI and API versions

###### Auto generated by spf13/cobra on 2-Aug-2019
