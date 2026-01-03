# Malang
### Yet another installer for arch linux

Malang is an experimental arch linux written in go for experimental and educational purposes.
The installer make use of 

- Go language
- Bubbletea TUI framework
- Linux commands
- Arch linux (to build iso)

## Installer Development

To build the installer, you need to have a working arch linux environment with go installed. Distrobox is recommended for easy setup.

There is an `Archlinux` dockerfile available to quickly setup a development environment. Build the docker image and then use it to create archlinux
distrobox container.

Once in the distrobox, you can use `/bin/build_iso` to build the iso.

There are other utilities in `installer/bin` to play around with the installer.

- `installer/bin/run` - Run the installer in the current environment (for testing, be careful with partitioning)
- `installer/bin/run_mock` - Run the installer in mock mode (no partitioning, no installation)
- `installer/bin/run_cli` - Run the installer in CLI mode to test the commands

## Features

- Network selection (wifi/ethernet)
- Disk partitioning (cfdisk)
- User setup

### Screenshots
<img width="1002" height="482" alt="Screenshot 2026-01-02 224050" src="https://github.com/user-attachments/assets/a16f7495-d916-48c1-9d6d-9c3ac5b316c1" />
<img width="787" height="517" alt="Screenshot 2026-01-02 224039" src="https://github.com/user-attachments/assets/b20c61b6-a8d2-494a-b583-f348b213f089" />
<img width="921" height="423" alt="Screenshot 2026-01-02 224016" src="https://github.com/user-attachments/assets/2f51673d-6156-4cde-9136-2f779fc1fca1" />
<img width="825" height="395" alt="Screenshot 2026-01-02 223957" src="https://github.com/user-attachments/assets/5df4f615-8d5a-4964-9c72-18fe9ab5c2d4" />
<img width="868" height="458" alt="Screenshot 2026-01-02 224101" src="https://github.com/user-attachments/assets/b6e5d7bb-a405-424d-91ce-c0c473d7aa5b" />
