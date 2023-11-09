package connections

import (
	"fmt"
	"github.com/aousomran/sqlite-og/internal/cbchannels"
	"github.com/aousomran/sqlite-og/internal/dbwrapper"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"strings"
	"sync"
)

type Manager struct {
	mutex  sync.RWMutex
	CnxMap map[string]*dbwrapper.DBWrapper
}

func NewManager() *Manager {
	return &Manager{
		mutex:  sync.RWMutex{},
		CnxMap: map[string]*dbwrapper.DBWrapper{},
	}
}

func (m *Manager) addConnection(id string, cnx *dbwrapper.DBWrapper) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.CnxMap[id]
	if ok {
		slog.Warn("connection already exists, replacing", "id", id)
	}
	m.CnxMap[id] = cnx
}

func (m *Manager) GetConnection(id string) (*dbwrapper.DBWrapper, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	cnx, ok := m.CnxMap[id]
	if !ok {
		return nil, fmt.Errorf("connection with id `%s` does not exist", id)
	}
	if cnx == nil {
		return nil, fmt.Errorf("dbwrapper in map is nil")
	}
	return cnx, nil
}

func (m *Manager) DeleteConnection(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.CnxMap, id)
}

func (m *Manager) Connect(dbname string, functions []string, aggregators []string) (string, error) {
	id := strings.Split(uuid.New().String(), "-")[0]
	channels := cbchannels.New()
	cnx := dbwrapper.New(dbname, id, functions, channels)
	err := cnx.Open(id)
	if err != nil {
		return "", err
	}
	m.addConnection(id, cnx)
	return id, nil
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var err error
	for id, cnx := range m.CnxMap {
		err = cnx.Close()
		if err != nil {
			slog.Error("cannot close connection", "error", err, "cnx_id", id, "dbname", cnx.Name)
		}
	}
	return err
}
