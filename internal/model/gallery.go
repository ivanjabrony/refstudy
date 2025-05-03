package model

type Picture struct {
	Name   string
	Path   string //url or on disk??
	Tags   string
	Height int
	Width  int
}

type Gallery struct {
	GalleryName string
	Description string
	IsPublic    bool
	Pictures    []Picture
	CurrentSize int
	OwnerName   string
}

type PictureTag struct {
	TagName string
}
