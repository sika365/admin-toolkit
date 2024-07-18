package registrar

type Package interface {
	Init() error
	Name() string
	Migrator() error
	Error() error
}

type Packages []Package
