package api

import (
	"errors"
	"sync"
)

type HashRepository struct {
	mtx       sync.RWMutex
	keyGen    *AtomicInt
	hashItems []HashItem
}

func NewHashRepository() *HashRepository {
	hr := &HashRepository{
		mtx:       sync.RWMutex{},
		keyGen:    NewAtomicInt(),
		hashItems: []HashItem{},
	}
	return hr
}

func (hr *HashRepository) Add(hi *HashItem) (id int64) {
	hi.ID = hr.keyGen.Inc()
	hr.mtx.Lock()
	hr.hashItems = append(hr.hashItems, *hi)
	hr.mtx.Unlock()
	return hi.ID
}

func (hr *HashRepository) Get(id int64) (hi *HashItem, err error) {
	if id < 1 || id > hr.Count() {
		return nil, errors.New("Index out of bounds")
	}
	hr.mtx.RLock()
	hi = &hr.hashItems[id-1]
	hr.mtx.RUnlock()
	return hi, nil
}

func (hr *HashRepository) Count() int64 {
	return hr.keyGen.Get()
}
