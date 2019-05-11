package spec

var BaseConfig = Spec{
	Root: Root{
		Path:     "/rootfs",
		ReadOnly: false,
	},
	Mounts: []Mount{
		Mount{
			Source:      "proc",
			Destination: "/proc",
			MountType:   "proc",
		},
		Mount{
			Source:      "tmpfs",
			Destination: "/dev",
			MountType:   "tmpfs",
			Options:     []string{"size=65536k", "mode=755"},
		},
		Mount{
			Source:      "/bin",
			Destination: "/bin",
			MountType:   "bind",
		},
		Mount{
			Source:      "/lib",
			Destination: "/lib",
			MountType:   "bind",
		},
		Mount{
			Source:      "/lib64",
			Destination: "/lib64",
			MountType:   "bind",
		},
		Mount{
			Source:      "/usr/lib",
			Destination: "/usr/lib",
			MountType:   "bind",
		},
	},
	Process: Process{
		CWD:  "/",
		Args: []string{"/bin/bash"},
		Env: []string{
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"TERM=xterm",
		},
	},
	Hostname: "freighter",
	Linux: Linux{
		Devices: []Device{
			Device{
				Path:     "/dev/null",
				DevType:  "c",
				Major:    1,
				Minor:    3,
				FileMode: 0666,
			},
			Device{
				Path:     "/dev/zero",
				DevType:  "c",
				Major:    1,
				Minor:    5,
				FileMode: 0666,
			},
			Device{
				Path:     "/dev/full",
				DevType:  "c",
				Major:    1,
				Minor:    3,
				FileMode: 0666,
			},
			Device{
				Path:     "/dev/random",
				DevType:  "c",
				Major:    1,
				Minor:    7,
				FileMode: 0666,
			},
			Device{
				Path:     "/dev/urandom",
				DevType:  "c",
				Major:    1,
				Minor:    8,
				FileMode: 0666,
			},
			Device{
				Path:     "/dev/tty",
				DevType:  "c",
				Major:    1,
				Minor:    9,
				FileMode: 0666,
			},
		},
		RootFSPropagation: "rslave",
	},
}
