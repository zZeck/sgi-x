# SGI Extractor

This is a fork of https://github.com/sgi-demos/sgix.

Extracts from an `*.sw` and `*.idb` set of files.

Multiple versions of sgix will be in this tree, one for each major IRIX release (3, 4, 5, 6).  

Expect hiccups, this is only complete in terms of getting it to work sufficiently to extract SGI demo source.

## Building

```
https://hexfiend.com/ is helpful for debugging (matching up .idb file info to .sw file reality).

brew install go
$ go get
$ go build
```

## Extracting Files

Let's say you have a `*.sw` and `*.idb` file. It's an extracted "tardist" file from SGI IRIX, or something like that. I am not an expert, but the whole system turns out to be fairly obvious to reverse engineer.

```
$ ls
Example
Example.idb
Example.sw
```

To extract, just run `sgix` and specify the destination:

```
$ sgix Example.idb Example.sw out
```

This will create a folder called `out` with the extracted contents.

## License

Licensed under the MIT license. See `LICENSE.txt`.

## See Also

http://persephone.cps.unizar.es/~spd/src/other/mydb.c
https://github.com/sgi-demos/sgix
