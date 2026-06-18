package categoryentity

const (
	initialID = 0
)

type Category struct {
	ID       int
	ParentID *int
	Name     string
}

func NewInitial(
	parentID *int,
	name string,
) Category {
	return Category{
		ID:       initialID,
		ParentID: parentID,
		Name:     name,
	}
}
