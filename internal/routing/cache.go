package routing

import (
	"sync"
	"sync/atomic"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

type cacheKey struct {
	From  model.SatelliteID
	To    model.SatelliteID
	Epoch uint64
}

type cacheEntry struct {
	key      cacheKey
	result   model.RouteResult
	prev     *cacheEntry
	next     *cacheEntry
	valid    uint32
}

type RouteCache struct {
	mu       sync.RWMutex
	capacity int
	items    map[cacheKey]*cacheEntry
	head     *cacheEntry
	tail     *cacheEntry

	currentEpoch uint64

	hits   atomic.Uint64
	misses atomic.Uint64
	evicts atomic.Uint64
}

func NewRouteCache(capacity int) *RouteCache {
	rc := &RouteCache{
		capacity: capacity,
		items:    make(map[cacheKey]*cacheEntry, capacity),
	}
	return rc
}

func (rc *RouteCache) SetEpoch(epoch uint64) {
	atomic.StoreUint64(&rc.currentEpoch, epoch)
}

func (rc *RouteCache) GetEpoch() uint64 {
	return atomic.LoadUint64(&rc.currentEpoch)
}

func (rc *RouteCache) Get(from, to model.SatelliteID, matrixEpoch uint64) (model.RouteResult, bool) {
	key := cacheKey{From: from, To: to, Epoch: matrixEpoch}

	rc.mu.RLock()
	entry, ok := rc.items[key]
	if !ok {
		rc.mu.RUnlock()
		rc.misses.Add(1)
		return model.RouteResult{}, false
	}

	if atomic.LoadUint32(&entry.valid) == 0 {
		rc.mu.RUnlock()
		rc.misses.Add(1)
		return model.RouteResult{}, false
	}

	if entry.key.Epoch != matrixEpoch {
		rc.mu.RUnlock()
		rc.misses.Add(1)
		return model.RouteResult{}, false
	}

	if entry.prev != nil {
		result := entry.result
		rc.mu.RUnlock()

		rc.mu.Lock()
		rc.moveToFront(entry)
		rc.mu.Unlock()

		rc.hits.Add(1)
		return result, true
	}

	result := entry.result
	rc.mu.RUnlock()
	rc.hits.Add(1)
	return result, true
}

func (rc *RouteCache) Put(from, to model.SatelliteID, matrixEpoch uint64, result model.RouteResult) {
	key := cacheKey{From: from, To: to, Epoch: matrixEpoch}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	if entry, ok := rc.items[key]; ok {
		entry.result = result
		atomic.StoreUint32(&entry.valid, 1)
		rc.moveToFront(entry)
		return
	}

	entry := &cacheEntry{
		key:    key,
		result: result,
		valid:  1,
	}

	rc.items[key] = entry
	rc.addToFront(entry)

	if len(rc.items) > rc.capacity {
		rc.evictTail()
	}
}

func (rc *RouteCache) InvalidateEpoch(epoch uint64) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	count := 0
	for key, entry := range rc.items {
		if key.Epoch == epoch {
			atomic.StoreUint32(&entry.valid, 0)
			delete(rc.items, key)
			rc.remove(entry)
			count++
		}
	}
	rc.evicts.Add(uint64(count))
}

func (rc *RouteCache) InvalidateOlderThan(epoch uint64) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	count := 0
	for key, entry := range rc.items {
		if key.Epoch < epoch {
			atomic.StoreUint32(&entry.valid, 0)
			delete(rc.items, key)
			rc.remove(entry)
			count++
		}
	}
	rc.evicts.Add(uint64(count))
}

func (rc *RouteCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	for _, entry := range rc.items {
		atomic.StoreUint32(&entry.valid, 0)
	}
	rc.items = make(map[cacheKey]*cacheEntry, rc.capacity)
	rc.head = nil
	rc.tail = nil
}

func (rc *RouteCache) Len() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return len(rc.items)
}

func (rc *RouteCache) Stats() (hits, misses, evicts uint64) {
	return rc.hits.Load(), rc.misses.Load(), rc.evicts.Load()
}

func (rc *RouteCache) addToFront(entry *cacheEntry) {
	entry.next = rc.head
	entry.prev = nil
	if rc.head != nil {
		rc.head.prev = entry
	}
	rc.head = entry
	if rc.tail == nil {
		rc.tail = entry
	}
}

func (rc *RouteCache) moveToFront(entry *cacheEntry) {
	if entry == rc.head {
		return
	}
	rc.remove(entry)
	rc.addToFront(entry)
}

func (rc *RouteCache) remove(entry *cacheEntry) {
	if entry.prev != nil {
		entry.prev.next = entry.next
	} else {
		rc.head = entry.next
	}
	if entry.next != nil {
		entry.next.prev = entry.prev
	} else {
		rc.tail = entry.prev
	}
	entry.prev = nil
	entry.next = nil
}

func (rc *RouteCache) evictTail() {
	if rc.tail == nil {
		return
	}
	entry := rc.tail
	atomic.StoreUint32(&entry.valid, 0)
	delete(rc.items, entry.key)
	rc.remove(entry)
	rc.evicts.Add(1)
}

type VersionedRouteResult struct {
	Result model.RouteResult
	Epoch  uint64
}
