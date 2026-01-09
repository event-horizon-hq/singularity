package manager

import (
	"os"
	"path/filepath"
	"singularity/internal/data"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type BlueprintManager struct {
	blueprintCache map[string]*data.Blueprint
	cacheMutex     sync.RWMutex
}

func CreateNewBlueprintManager() *BlueprintManager {
	return &BlueprintManager{
		blueprintCache: make(map[string]*data.Blueprint),
		cacheMutex:     sync.RWMutex{},
	}
}

func (blueprintManager *BlueprintManager) GetBlueprint(id string) (*data.Blueprint, bool) {
	blueprintManager.cacheMutex.RLock()
	defer blueprintManager.cacheMutex.RUnlock()

	blueprint, ok := blueprintManager.blueprintCache[id]

	return blueprint, ok
}

func (blueprintManager *BlueprintManager) GetAllBlueprints() []*data.Blueprint {
	blueprintManager.cacheMutex.RLock()
	defer blueprintManager.cacheMutex.RUnlock()

	blueprints := make([]*data.Blueprint, 0, len(blueprintManager.blueprintCache))
	for _, bp := range blueprintManager.blueprintCache {
		blueprints = append(blueprints, bp)
	}

	return blueprints
}

func (blueprintManager *BlueprintManager) LoadBlueprint(path string) (*data.Blueprint, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var blueprint data.Blueprint

	if err := toml.Unmarshal(file, &blueprint); err != nil {
		return nil, err
	}

	return &blueprint, nil
}

func (blueprintManager *BlueprintManager) LoadBlueprints(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".toml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		blueprint, err := blueprintManager.LoadBlueprint(path)
		if err != nil {
			return false, err
		}

		blueprintManager.cacheMutex.Lock()
		blueprintManager.blueprintCache[blueprint.Id] = blueprint
		blueprintManager.cacheMutex.Unlock()
	}

	return true, nil
}
