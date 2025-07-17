package command

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sync"
)

var (
	ErrCreateApplicationCommand = errors.New("failed to create discord application command")
	ErrDeleteApplicationCommand = errors.New("failed to delete discord application command")
)

type Registry struct {
	mutex    sync.RWMutex
	commands map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

func (r *Registry) Register(session *discordgo.Session, command Command) error {
	discordCmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.ApplicationCommand())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateApplicationCommand, err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.commands[discordCmd.ID] = command

	return nil
}

func (r *Registry) Unregister(session *discordgo.Session, commandID string) error {
	err := session.ApplicationCommandDelete(session.State.User.ID, "", commandID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDeleteApplicationCommand, err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.commands, commandID)

	return nil
}

func (r *Registry) Commands() map[string]Command {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	cmds := make(map[string]Command)
	for k, v := range r.commands {
		cmds[k] = v
	}

	return cmds
}
