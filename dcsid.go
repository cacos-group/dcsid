package dcsid

import (
	"context"
	"fmt"
	"github.com/cacos-group/dcsid/sql"
	"log"
	"sync"
)

type idPool struct {
	NextID uint64
	MaxID  uint64
	Step   int
}

type Client struct {
	mu sync.Mutex

	idPool0 *idPool
	idPool1 *idPool

	idPool1Chan chan bool
	idPool1Mu   sync.Mutex

	mysql *sql.MySQL

	config *Config
}

type Config struct {
	DSN    string // write data source name.
	BizTag string
	Step   int
}

func New(config *Config) *Client {
	ch := make(chan bool, 1)

	ms := sql.NewMySQL(&sql.Config{
		DSN:    config.DSN,
		BizTag: config.BizTag,
	})

	if config.Step <= 100 {
		config.Step = 1000
	}

	err := ms.InitBizTag(context.TODO(), config.BizTag, 1, config.Step, "")
	if err != nil {
		log.Panicf("InitBizTag err: %v\n", err)
	}

	l := &Client{
		config:      config,
		idPool1Chan: ch,
		mysql:       ms,
	}

	idPool, err := l.generateIdPool()
	if err != nil {
		log.Panicf("initIdPool0 err: %v\n", err)
	}

	l.idPool0 = idPool

	idPool1, err := l.generateIdPool()
	if err != nil {
		log.Panicf("checkIdPool1 err: %v\n", err)
	}

	l.idPool1 = idPool1

	l.initDaemon()

	return l
}

func (l *Client) NextId() (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.getNextId()
}

var IDPool0EmptyErr = fmt.Errorf("idPool0 empty")

func (l *Client) getNextId() (uint64, error) {
	if l.idPool0 != nil {
		nextID := l.idPool0.NextID
		if nextID <= l.idPool0.MaxID {
			l.idPool0.NextID = nextID + 1
			return nextID, nil
		}
	}

	l.checkIdPool1()

	l.idPool0 = l.idPool1
	l.idPool1 = nil
	// 通知生成idPool1
	l.idPool1Chan <- true

	nextID := l.idPool0.NextID
	if nextID <= l.idPool0.MaxID {
		l.idPool0.NextID = nextID + 1
		return nextID, nil
	}

	return 0, IDPool0EmptyErr
}

func (l *Client) initDaemon() {
	go func() {
		for {
			select {
			case <-l.idPool1Chan:
				l.checkIdPool1()
			}
		}
	}()
}

func (l *Client) checkIdPool1() {
	if l.idPool1 == nil {
		// 同步生成idPool1
		l.idPool1Mu.Lock()
		if l.idPool1 == nil {
			idPool, err := l.generateIdPool()
			if err != nil {
				log.Panicf("checkIdPool1 err: %v\n", err)
			}

			l.idPool1 = idPool
		}
		l.idPool1Mu.Unlock()
	}
}

func (l *Client) generateIdPool() (*idPool, error) {
	startID, endID, step, err := l.mysql.GetEndID(context.TODO())
	if err != nil {
		return nil, err
	}

	maxID := endID
	currentID := startID

	return &idPool{
		NextID: currentID,
		MaxID:  maxID,
		Step:   step,
	}, nil
}
