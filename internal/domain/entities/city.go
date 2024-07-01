package entities

type City struct {
	ID    int
	Name  string `xml:"nome"`
	State string
}
