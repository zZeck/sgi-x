# SGI Extractor

This is a fork of https://github.com/depp/sgix.

Extracts from `*.sw`, `*.idb`, etc. sets of SGI install files. 

This fork aims to expand handling of install files to all major IRIX releases (3, 4, 5, 6).  

The goal is to be able to extract SGI demo source and man pages, so it has only been tested on install files containing those things.

## Building

A hex editor (like https://hexfiend.com/) is useful for debugging (getting the info from the .idb to match the reality of the .sw and .man files).

```
$ brew install go
$ go get
$ go build
```

## Extracting Files

Let's say you have a `*.sw` and `*.idb` file. It's an extracted "tardist" file from SGI IRIX, or something like that. I am not an expert, but the whole system turns out to be fairly obvious to reverse engineer.

```
$ ls
dev.idb
dev.sw
```

To extract, just run `sgix` and specify the destination:

```
$ sgix dev.idb dev.sw out
```

This will create a folder called `out` with the extracted contents.

## License

Licensed under the MIT license. See `LICENSE.txt`.

## See Also

http://persephone.cps.unizar.es/~spd/src/other/mydb.c
https://github.com/depp/sgix
