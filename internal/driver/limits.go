package driver

import (
	"strings"

	"github.com/linode/linode-blockstorage-csi-driver/pkg/hwinfo"
)

// maxVolumeAttachments returns the maximum number of block storage volumes
// that can be attached to a Linode instance, given the amount of memory the
// instance has.
//
// TODO: This code should be cleaned up to use the built-in max and min
// functions once the project is updated to Go 1.21. See
// https://go.dev/ref/spec#Min_and_max.
func maxVolumeAttachments(memoryBytes uint) int {
	attachments := memoryBytes >> 30
	if attachments > maxAttachments {
		return maxAttachments
	}
	if attachments < maxPersistentAttachments {
		return maxPersistentAttachments
	}
	return int(attachments)
}

const (
	// maxPersistentAttachments is the default number of volume attachments
	// allowed when they are persisted to an instance/boot config. This is
	// also the maximum number of allowed volume attachments when the
	// instance type has < 16GiB of RAM.
	maxPersistentAttachments = 8

	// maxAttachments it the hard limit of volumes that can be attached to
	// a single Linode instance.
	maxAttachments = 64
)

func attachedVolumeCount(hw hwinfo.HardwareInfo) (int, error) {
	bdev, err := hw.Block()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, disk := range bdev.Disks {
		driveType := strings.ToLower(disk.DriveType.String())
		// Skip loopbacks, virtual devices, CD-ROMs, and removable drives
		if disk.IsRemovable || driveType == "virtual" || driveType == "cdrom" {
			continue
		}

		// Only consider disks that have partitions and aren't read-only
		if len(disk.Partitions) > 0 {
			count++
		}

	}
	return count, nil
}
