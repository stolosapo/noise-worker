package example

type (
	TheModel interface {
		ID() int
		Name() string
		Description() string

		Value() int
		WithValue(value int) TheModel
	}

	theModel struct {
		id          int
		name        string
		description string
		value       int
	}
)

func NewModel(
	id int,
	name string,
	description string,
) *theModel {
	return &theModel{
		id:          id,
		name:        name,
		description: description,
		value:       0,
	}
}

func (m theModel) ID() int {
	return m.id
}

func (m theModel) Name() string {
	return m.name
}

func (m theModel) Description() string {
	return m.description
}

func (m theModel) Value() int {
	return m.value
}

func (m theModel) WithValue(value int) TheModel {
	newModel := NewModel(
		m.id,
		m.name,
		m.description,
	)
	newModel.value = value
	return newModel
}
