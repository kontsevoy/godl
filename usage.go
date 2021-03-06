package main

const UsageStr = `
godl: quick&dirty ACI builder with ability to add dependencies (so-libs) to rootfs
If the ACI file already exists, it appends to it unless -f is given.

USAGE:
    godl [OPTION]... PATTERN PATTERN...

PATTERN
    File(s) or directories to pack into the ACI

OPTIONS:
    -r : output directory which will contain rootfs + manifest (./aci)
    -o : output ACI image (none)
    -t : target directory within rootfs (/)
    -f : force overwrite. (false)
    -m : manifest file to use (auto-generates one if missing)
    -n : application name, if no manifest given
    -n : application description, if no manifest given
    -i : ignore binary dependencies (false)

EXAMPLES:
    godl -o ed.aci /bin/ed
    godl -r dir ../project
`
