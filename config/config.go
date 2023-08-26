package config

import (
	"os"
)

type Conf struct {
	RMQ       *RMQ
	Image     *Image
	MountPath string
}

func NewConf() *Conf {
	return &Conf{
		RMQ:       newRMQ(),
		Image:     newImage(),
		MountPath: os.Getenv("MOUNT_PATH"),
	}
}
