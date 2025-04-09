# pds

A Project Directory Generator written in go

To build the project, run the following
```sh
make build
```

The command will generate a binary file `./bin/pds`.

You can then move the binary to your local path to use it anywhere.

For the directory structure, you'll need a `structure.toml` file that defines the directory structure you want.

Let's create one with a dummy content below to see how it it creates the directory structure.
```toml
files = [
    "README.md",
    "Makefile"
]

[bin]

[src]
files = [ "main.go" ]

[test]
[test.data]
```

Once you save the content above in `structure.toml` file, you can run the following command in the same directory as of the `structure.toml` file:

```sh
pds 
```

This will create a directory called `new_project` and will have the following contents:
```sh
|- bin
|- src
    |- main.go
|- test
    |- data
|- Makefile
|- README.md
```

This demonstrates how you can create the files in the base project directory, empty directories, files in a directory, and nested empty directories respectively.
