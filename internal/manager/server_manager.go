package manager

import (
	"context"
	"singularity/internal/data"
	"singularity/internal/enum"
	"singularity/internal/repository"
	"sync"
)

type ServerManager struct {
	cache      map[string]*data.Server
	cacheMutex sync.RWMutex
	repository *repository.ServerRepository
}

func CreateNewServerManager(repository *repository.ServerRepository) *ServerManager {
	return &ServerManager{
		cache:      make(map[string]*data.Server),
		cacheMutex: sync.RWMutex{},
		repository: repository,
	}
}

func (serverManager *ServerManager) GetServer(serverId string) (*data.Server, error) {
	serverManager.cacheMutex.RLock()
	defer serverManager.cacheMutex.RUnlock()

	server, ok := serverManager.cache[serverId]
	if !ok {
		return nil, nil
	}

	return server, nil
}

func (serverManager *ServerManager) GetAllServers() []*data.Server {
	serverManager.cacheMutex.RLock()
	defer serverManager.cacheMutex.RUnlock()

	servers := make([]*data.Server, 0, len(serverManager.cache))
	for _, server := range serverManager.cache {
		servers = append(servers, server)
	}

	return servers
}

func (serverManager *ServerManager) AddServer(server *data.Server) bool {
	serverManager.cacheMutex.Lock()
	defer serverManager.cacheMutex.Unlock()

	id := server.Id()
	if _, exists := serverManager.cache[id]; exists {
		return false
	}

	err := serverManager.repository.Insert(context.Background(), server)
	if err != nil {
		return false
	}

	serverManager.cache[id] = server
	return true
}

func (serverManager *ServerManager) LoadServer(server *data.Server) bool {
	serverManager.cacheMutex.Lock()
	defer serverManager.cacheMutex.Unlock()

	id := server.Id()
	if _, exists := serverManager.cache[id]; exists {
		return false
	}

	serverManager.cache[id] = server
	return true
}

func (serverManager *ServerManager) DeleteServer(id string) (*data.Server, bool) {
	serverManager.cacheMutex.Lock()
	defer serverManager.cacheMutex.Unlock()

	server, exists := serverManager.cache[id]
	if !exists {
		return nil, false
	}

	err := serverManager.repository.DeleteByID(context.Background(), id)
	if err != nil {
		return nil, false
	}

	delete(serverManager.cache, id)
	return server, true
}

func (serverManager *ServerManager) UpdateStatus(id string, status enum.Status) bool {
	serverManager.cacheMutex.Lock()
	defer serverManager.cacheMutex.Unlock()

	s, ok := serverManager.cache[id]
	if !ok {
		return false
	}
	s.Status = status
	return true
}

func (serverManager *ServerManager) UpdateReport(id string, report data.ServerReport) bool {
	serverManager.cacheMutex.Lock()
	defer serverManager.cacheMutex.Unlock()

	s, ok := serverManager.cache[id]
	if !ok {
		return false
	}
	s.Report = &report
	return true
}