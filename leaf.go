package leaf

import (
	"context"
	"fmt"
	"github.com/go-light/leaf/sql"
	"log"
	"sync"
)

type idPool struct {
	NextID uint64
	MaxID  uint64
	Step   int
}

type Leaf struct {
	mu           sync.Mutex
	idPool1Mutex sync.Mutex

	idPool0 *idPool
	idPool1 *idPool
	ch      chan bool

	mysql *sql.MySQL

	config *Config
}

type Config struct {
	DSN    string // write data source name.
	BizTag string
}

func New(config *Config) *Leaf {
	ch := make(chan bool, 2000)

	ms := sql.NewMySQL(&sql.Config{
		DSN:    config.DSN,
		BizTag: config.BizTag,
	})

	leaf := &Leaf{
		config: config,
		ch:     ch,
		mysql:  ms,
	}

	err := leaf.initIdPool0()
	if err != nil {
		panic(err)
	}

	leaf.initDaemon()

	return leaf
}

func (l *Leaf) NextId() (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.getNextId()
}

var IDPool0EmptyErr = fmt.Errorf("idPool0 empty")

func (l *Leaf) getNextId() (uint64, error) {
	idPool := l.idPool0
	if idPool == nil {
		return 0, IDPool0EmptyErr
	}

	nextID := idPool.NextID
	if nextID <= idPool.MaxID {
		l.idPool0.NextID = nextID + 1
		l.ch <- true
		return nextID, nil
	}

	if l.idPool1 == nil {
		err := l.initIdPool1()
		if err != nil {
			return 0, err
		}
	}

	l.idPool0 = l.idPool1
	l.idPool1 = nil

	return l.getNextId()
}

func (l *Leaf) initDaemon() {
	go func() {
		for {
			select {
			case <-l.ch:
				err := l.initIdPool1()
				if err != nil {
					log.Printf("daemon initIdPool1 err: %v", err)
				}
			}
		}
	}()
}

func (l *Leaf) initIdPool0() error {
	startID, endID, step, err := l.mysql.GetEndID(context.TODO())
	if err != nil {
		return err
	}

	maxID := endID
	currentID := startID

	l.idPool0 = &idPool{
		NextID: currentID,
		MaxID:  maxID,
		Step:   step,
	}

	return nil
}

func (l *Leaf) initIdPool1() error {
	l.idPool1Mutex.Lock()
	defer l.idPool1Mutex.Unlock()

	if l.idPool1 != nil && l.idPool1.MaxID > 0 {
		return nil
	}

	startID, endID, step, err := l.mysql.GetEndID(context.TODO())
	if err != nil {
		return err
	}

	maxID := endID
	currentID := startID

	l.idPool1 = &idPool{
		NextID: currentID,
		MaxID:  maxID,
		Step:   step,
	}

	return nil
}
