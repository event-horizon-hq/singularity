package manager

import (
	"context"
	"os"
	"path/filepath"
	"singularity/gen/blueprint"
	"strings"
	"sync"
)

type BlueprintManager struct {
	blueprintCache map[string]*blueprint.Blueprint
	cacheMutex     sync.RWMutex
}

func CreateNewBlueprintManager() *BlueprintManager {
	return &BlueprintManager{
		blueprintCache: make(map[string]*blueprint.Blueprint),
		cacheMutex:     sync.RWMutex{},
	}
}

func (blueprintManager *BlueprintManager) GetBlueprint(id string) (*blueprint.Blueprint, bool) {
	blueprintManager.cacheMutex.RLock()
	defer blueprintManager.cacheMutex.RUnlock()

	blueprint, ok := blueprintManager.blueprintCache[id]

	return blueprint, ok
}

func (blueprintManager *BlueprintManager) GetAllBlueprints() []*blueprint.Blueprint {
	blueprintManager.cacheMutex.RLock()
	defer blueprintManager.cacheMutex.RUnlock()

	blueprints := make([]*blueprint.Blueprint, 0, len(blueprintManager.blueprintCache))
	for _, bp := range blueprintManager.blueprintCache {
		blueprints = append(blueprints, bp)
	}

	return blueprints
}

func (blueprintManager *BlueprintManager) LoadBlueprint(path string) (*blueprint.Blueprint, error) {
	ctx := context.Background()

	bp, err := blueprint.LoadFromPath(ctx, path)
	if err != nil {
		return nil, err
	}

	return &bp, nil
}

func (blueprintManager *BlueprintManager) ReloadBlueprints(dir string) (int, error) {
	blueprintManager.cacheMutex.Lock()
	blueprintManager.blueprintCache = make(map[string]*blueprint.Blueprint)
	blueprintManager.cacheMutex.Unlock()

	_, err := blueprintManager.LoadBlueprints(dir)
	if err != nil {
		return 0, err
	}

	blueprintManager.cacheMutex.RLock()
	count := len(blueprintManager.blueprintCache)
	blueprintManager.cacheMutex.RUnlock()

	return count, nil
}

func (blueprintManager *BlueprintManager) LoadBlueprints(dir string) (bool, error) {
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if filepath.Ext(entry.Name()) != ".pkl" {
			return nil
		}

		// Skip base.pkl files (they are templates, not concrete blueprints)
		if strings.HasSuffix(entry.Name(), "base.pkl") || strings.HasPrefix(entry.Name(), "Reference") {
			return nil
		}

		bp, err := blueprintManager.LoadBlueprint(path)
		if err != nil {
			return err
		}

		blueprintManager.cacheMutex.Lock()
		blueprintManager.blueprintCache[bp.Id] = bp
		blueprintManager.cacheMutex.Unlock()

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
