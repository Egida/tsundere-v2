package roles

import (
	"log"
	"tsundere/source"
)

var (
	List = map[string]*Role{
		"admin": {
			DisplayName: "A",
			Color:       "#ffffff",
		},
	}
)

type Roles struct{}

func (r *Roles) Serve() {
	err := source.Parser.Unmarshal(&List, "roles")
	if err != nil {
		log.Fatal(err)
	}
}
