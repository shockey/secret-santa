# secret-santa 🤫🧑‍🎄

Use the power of depth-first search to match multiple groups of gift-givers among each other, with support for flexible matching rules ("exceptions").

### How to use

Assuming you have a Go environment set up, it's as easy as plugging your data into `config.go`, then running the program:

```
$ go run .
```

You'll find a new output CSV in the same directory as this README.

##### Real mode

CSVs will, by default, be written with `TEST` in their filenames. If you'd like them to be written as `REAL` instead of `TEST`, pass the `--real` CLI flag:

```
$ go run . --real
```

This is useful if you're giving the program a few dry runs while configuring your groups and exceptions.