package config

import "os"

func newImage() *Image {
	return &Image{
		Width:  os.Getenv("IMAGE_WIDTH"),
		Height: os.Getenv("IMAGE_HEIGHT"),
	}
}

type Image struct {
	Width  string
	Height string
}

func (c *Image) ImageInfo() *Image {
	return &Image{
		Width:  c.Width,
		Height: c.Height,
	}
}
