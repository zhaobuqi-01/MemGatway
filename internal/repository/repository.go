package repository

type Creator interface {
	Create(item any) error
}

type Updater interface {
	Update(item any) error
}

type Deleter interface {
	Delete(item any) error
}

type Getter interface {
	// Get(item T) (T, error)
	GetAll(item any) ([]any, error)
}

type Lister interface {
	List() ([]any, error)
}

type Tabler interface {
	TablerName() string
}
