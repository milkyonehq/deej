package player

import "sync"

type Registry struct {
	mutex   sync.RWMutex
	players map[string]Player
}

func NewRegistry() *Registry {
	return &Registry{
		players: make(map[string]Player),
	}
}

func (r *Registry) Unregister(guildID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.players[guildID].Stop()
	delete(r.players, guildID)
}

func (r *Registry) FindOrCreate(guildID string, createPlayer func() Player) Player {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if player, ok := r.players[guildID]; ok {
		return player
	}

	player := createPlayer()
	r.players[guildID] = player

	return player
}

func (r *Registry) Players() map[string]Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	players := make(map[string]Player)
	for k, v := range r.players {
		players[k] = v
	}

	return players
}
