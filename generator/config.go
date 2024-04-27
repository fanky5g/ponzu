package generator

type Path struct {
	Root string
	Base string
}

type Target struct {
	Path    Path
	Package string
}

type Config struct {
	Target Target
	Type   Type
}
