// +build cgo

package qemu

import (
	"os"
	"strconv"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/solo-io/unik/pkg/types"
)

func (p *QemuProvider) StopInstance(id string) error {
	instance, err := p.GetInstance(id)
	if err != nil {
		return errors.New("retrieving instance "+id, err)
	}

	// kill qemu
	pid, err := strconv.Atoi(instance.Id)
	if err != nil {
		return errors.New("invalid instance id (should be qemu pid)", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logrus.Warn("failed finding instance, assuming instance has externally terminated", err)
	} else {
		if err := process.Signal(syscall.SIGKILL); err != nil {
			logrus.Warn("failed terminating instance, assuming instance has externally terminated", err)
		}
	}

	volumesToDetach := []*types.Volume{}
	volumes, err := p.ListVolumes()
	if err != nil {
		return errors.New("getting volume list", err)
	}
	for _, volume := range volumes {
		if volume.Attachment == instance.Id {
			volumesToDetach = append(volumesToDetach, volume)
		}
	}

	return p.state.RemoveInstance(instance)
}
