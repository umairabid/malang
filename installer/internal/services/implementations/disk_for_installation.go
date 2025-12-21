package implementations

import (
    "os/exec"
    "strconv"
    "strings"

    types "installer.malang/internal/types"
)

func listDevices() []string {
    cmd := exec.Command("lsblk", "-b", "-d", "-o", "NAME,SIZE,TYPE,MOUNTPOINTS")
    out, err := cmd.Output()
    if err != nil {
        panic(err)
    }

    return strings.Split(string(out), "\n")
}

func parseDisks(lines []string) []types.Disk {
    var disks []types.Disk
    for _, line := range lines[1:] {
        fields := strings.Fields(line)
        if len(fields) < 2 {
            continue
        }

        name := fields[0]
        size, _ := strconv.ParseUint(fields[1], 10, 64)
        sizeInGb := size / (1024 * 1024 * 1024)
        formattedSizeInGb := strconv.FormatUint(sizeInGb, 10) + "G"

        dtype := fields[2]
        var mountPoint string

        if 3 < len(fields) {
            mountPoint = fields[3]
        }

        if sizeInGb > 0 {
            disks = append(disks, types.Disk{Name: name, Size: size, Type: dtype, MountPoint: mountPoint, SizeInGb: formattedSizeInGb})
        }
    }

    return disks
}

func CollectDisks() []types.Disk {
    lines := listDevices()
    return parseDisks(lines)
}
