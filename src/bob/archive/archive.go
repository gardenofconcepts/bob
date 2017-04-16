package archive

type Archive struct {
	path string
}

func NewArchive(path string) *Archive {
	return &Archive{
		path: path,
	}
}
