# Godl

Godl is a command line tool to quickly create ACI container images out of random data. 
Here are a few reasons why I wrote it:

##### Quick Binaries

Godl can be used to quickly package a given binary file on your host into an ACI.
It's quick because Godl:

1. Generates a passable Manifest, although you can specify your own.
2. Discover all dynamic libraries (dependencies) similar to `ldd` and will package them as well.
3. Takes random filename patterns, so you can quickly package entire subdirectories into ACI images.

##### Complimentary to actool

CoreOS folks have a standard ACI tool, which requires you to manually create rootfs and manifest. 
Godl can be complimentary to it: use godl to quickly generate RootFS content, inspect it and use actool 
to build an ACI. 

##### Limitations

1. Godl does not generate ACI signatures.
2. Godl is an awufl name

##### Usage

```
USAGE:
    godl [OPTION]... PATTERN 

PATTERN:
    Files and directories to pack into the ACI

OPTIONS:
    -r : output directory which will contain rootfs + manifest (./aci)
    -o : output ACI image (none)
    -t : target directory within rootfs. (/)
    -f : force overwrite. (false)
    -m : manifest file to use (auto-generates one if missing)
    -i : ignore binary dependencies (false)

EXAMPLES:
    godl -o ed.aci /bin/ed
    godl -r dir ../project/**/*
```

It's easy to play with it, using -r option (when it genertes a rootfs directory instead of an ACI) and 
inspecting it using `tree`. Also try `-t` option: it allows you to package a bunch of random files under a given
target dir within rootfs.

