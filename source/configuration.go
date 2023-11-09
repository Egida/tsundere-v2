package source

import "tsundere/packages/utilities/simpleconfig"

const (
	NAME    string = "Tsundere"
	VERSION string = "v1.0.0-dev"
	HEADER  string = `   ______                     __             
  /_  __/______  ______  ____/ /__  ________ 
   / / / ___/ / / / __ \/ __  / _ \/ ___/ _ \
  / / (__  ) /_/ / / / / /_/ /  __/ /  /  __/
 /_/ /____/\__,_/_/ /_/\__,_/\___/_/   \___/  ` + VERSION + `

`

	Build = DEVELOPER
)

const (
	RELEASE = iota
	BETA
	PRE_ALPHA
	DEVELOPER
)

var (
	Parser = simpleconfig.New(simpleconfig.Json, "resources")
)
