FROM archlinux:latest

# Install necessary packages
RUN pacman -Syu --noconfirm archiso git base-devel rsync go parted
