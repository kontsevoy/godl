package main

const UsageStr = `
godl: quick&dirty ACI builder with ability to add dependencies (so-libs) to rootfs
If the ACI file already exists, it appends to it unless -f is given.

USAGE:
    godl [OPTION]... PATTERN 

PATTERN:
    Files and directories to pack into the ACI

OPTIONS:
    -o : output (./out.aci)
    -r : output directory, what out.aci would contain (./aci)
    -t : target directory within rootfs. (/)
    -f : force overwrite. (false)
    -m : manifest file to use (auto-generates one if missing)

EXAMPLES:
    godl -o ed.aci /bin/ed
    godl -r dir ../project/**/*
`
