package main

const (
	// DefaultManifest will be used to create manifest file if
	// user doesn't specify one
	DefaultManifest = `{
    "acKind": "ImageManifest",
    "acVersion": "0.5.2",
    "annotations": null,
    "app": {
        "environment": [],
        "eventHandlers": null,
        "exec": [
            "/bin/plugin"
        ],
        "group": "0",
        "isolators": null,
        "mountPoints": null,
        "ports": null,
        "user": "0"
    },
    "dependencies": null,
    "labels": [
        {
            "name": "os",
            "value": "linux"
        },
        {
            "name": "arch",
            "value": "amd64"
        }
    ],
    "name": "%app-name%",
    "pathWhitelist": null
}
`
)
