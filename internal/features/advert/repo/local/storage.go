package advertlocal

type Storage struct {
	root string
}

func New(root string) *Storage {
	return &Storage{
		root: root,
	}
}
