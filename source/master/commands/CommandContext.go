package commands

import "tsundere/source/database"

type CommandContext struct {
	arguments map[string]*ParsedArgument
}

func (ctx *CommandContext) User(name string) (*database.UserProfile, error) {
	return nil, nil
}
