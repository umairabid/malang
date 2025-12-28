# Malang
### Yet another installer for arch linux

Malang is an experimental arch linux written in go for experimental and educational purposes.
The installer make use of 

- Go language
- Bubbletea TUI framework
- Linux commands
- Arch linux (to build iso)

## Features

- Network selection (wifi/ethernet)
- Disk partitioning (cfdisk)
- User setup

## Installer Development

To build the installer, you need to have a working arch linux environment with go installed. Distrobox is recommended for easy setup.

There is an `Archlinux` dockerfile available to quickly setup a development environment. Build the docker image and then use it to create archlinux
distrobox container.

Once in the distrobox, you can use `/bin/build_iso` to build the iso.

There are other utilities in `installer/bin` to play around with the installer.

- `installer/bin/run` - Run the installer in the current environment (for testing, be careful with partitioning)
- `installer/bin/run_mock` - Run the installer in mock mode (no partitioning, no installation)
- `installer/bin/run_cli` - Run the installer in CLI mode to test the commands
