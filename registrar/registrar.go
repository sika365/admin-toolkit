package registrar

import (
	"errors"
	"strings"
)

var (
	ErrPackageAlreadyExists = errors.New("package already exists")
	ErrNotFound             = errors.New("package not found")
	ErrInvalidName          = errors.New("invalid package name")
)

type Registrar interface {
	Init() error
	Add(pkgs ...Package) Registrar
	Replace(name string, pkg Package) Registrar
	Del(name string) Registrar
	Get(name string) (Package, error)
	Error() error
}

type reg struct {
	packages map[string]Package
	err      error
}

var registrar *reg

func init() {
	registrar = New().(*reg)
}

func New() Registrar {
	r := &reg{
		packages: make(map[string]Package),
		err:      nil,
	}
	return r
}

func Instance() *reg {
	return registrar
}

func Error() error { return registrar.Error() }
func (r *reg) Error() error {
	return r.err
}

func Init() (err error) { return registrar.Init() }
func (r *reg) Init() (err error) {
	for _, p := range r.packages {
		if err := p.Init(); err != nil {
			return err
		}
	}
	return nil
}

// Get returns package by name
func Get(name string) (pkg Package, err error) { return registrar.Get(name) }

// Get returns package by name
func (r *reg) Get(name string) (pkg Package, err error) {
	if r.err != nil {
		return nil, r.err
	} else if name = strings.TrimSpace(name); name == "" {
		return nil, ErrInvalidName
	} else if pkg, ok := r.packages[name]; !ok {
		return nil, ErrNotFound
	} else {
		return pkg, nil
	}
}

// Add adds package in registrar
func Add(pkgs ...Package) Registrar { return registrar.Add(pkgs...) }

// Add adds package in registrar
func (r *reg) Add(pkgs ...Package) Registrar {
	if r.err != nil {
		return r
	} else {
		for _, pkg := range pkgs {
			if _, r.err = r.Get(pkg.Name()); r.err != nil && !errors.Is(r.err, ErrNotFound) {
				return r
			} else if errors.Is(r.err, ErrNotFound) {
				r.packages[pkg.Name()] = pkg
			}
		}
		r.err = nil
		return r
	}
}

// Del implements Registrar.
func Del(name string) Registrar { return registrar.Del(name) }

// Del implements Registrar.
func (r *reg) Del(name string) Registrar {
	if r.err != nil {
		return r
	} else if _, r.err = r.Get(name); r.err != nil {
		return r
	} else {
		delete(r.packages, name)
		return r
	}
}

// Replace implements Registrar.
func Replace(name string, newPkg Package) Registrar { return registrar.Replace(name, newPkg) }

// Replace implements Registrar.
func (r *reg) Replace(name string, newPkg Package) Registrar {
	if err := r.Del(name).Error(); err != nil && !errors.Is(err, ErrNotFound) {
		return r
	} else if err := r.Add(newPkg).Error(); err != nil {
		return r
	} else {
		return r
	}
}
