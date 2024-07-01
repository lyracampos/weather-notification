package entities

type Weather struct {
	Day       string
	Condition string // TODO: criar um enum para exibir a descrição refernte a sigla
	Max       int
	Min       int
	IUV       float32
}
