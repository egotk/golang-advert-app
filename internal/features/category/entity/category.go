package categoryentity

const (
	initialID = 0
)

type Category struct {
	ID       int64
	ParentID *int64
	Name     string
}

func NewInitial(
	parentID *int64,
	name string,
) Category {
	return Category{
		ID:       initialID,
		ParentID: parentID,
		Name:     name,
	}
}
