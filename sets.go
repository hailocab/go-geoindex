package geoindex

import "time"

// A set interface.
type set interface {
	Add(id string)
	Has(id string) (ok bool)
	Remove(id string)
	IDs() []string
	Values() map[string]struct{}
	Size() int
	Clone() set
}

// A set that contains values.
type basicSet map[string]struct{}

func newSet() set {
	return basicSet(make(map[string]struct{}))
}

// Clone creates a copy of the set where the values in clone set point to the same underlying reference as the original set
func (set basicSet) Clone() set {
	clone := basicSet(make(map[string]struct{}, len(set)))
	for k, v := range set {
		clone[k] = v
	}

	return clone
}

func (set basicSet) Add(id string) {
	set[id] = struct{}{}
}

func (set basicSet) Remove(id string) {
	delete(set, id)
}

func (set basicSet) IDs() []string {
	result := make([]string, len(set))

	i := 0
	for key := range set {
		result[i] = key
		i++
	}

	return result
}

func (set basicSet) Values() map[string]struct{} {
	return set
}

func (set basicSet) Has(id string) (ok bool) {
	_, ok = set[id]
	return
}

func (set basicSet) Size() int {
	return len(set)
}

// An expiring set that removes the points after X minutes.
type expiringSet struct {
	values         set
	insertionOrder *queue
	expiration     Minutes
	onExpire       func(id string)
	lastInserted   map[string]time.Time
}

type timestampedValue struct {
	id        string
	timestamp time.Time
}

// Clone panics - We currently do not allow cloning of an expiry set
func (set *expiringSet) Clone() set {
	panic("Cannot clone an expiry set")
}

func newExpiringSet(expiration Minutes) *expiringSet {
	return &expiringSet{newSet(), newQueue(1), expiration, nil, make(map[string]time.Time)}
}

func (set *expiringSet) hasExpired(time time.Time) bool {
	currentTime := getNow()
	return int(currentTime.Sub(time).Minutes()) > int(set.expiration)
}

func (set *expiringSet) expire() {
	for !set.insertionOrder.IsEmpty() {
		lastInserted := set.insertionOrder.Peek().(*timestampedValue)

		if set.hasExpired(lastInserted.timestamp) {
			set.insertionOrder.Pop()

			if set.hasExpired(set.lastInserted[lastInserted.id]) {
				set.values.Remove(lastInserted.id)

				if set.onExpire != nil {
					set.onExpire(lastInserted.id)
				}
			}
		} else {
			break
		}
	}
}

func (set *expiringSet) Add(id string) {
	set.expire()
	set.values.Add(id)
	insertionTime := getNow()
	set.lastInserted[id] = insertionTime
	set.insertionOrder.Push(&timestampedValue{id, insertionTime})
}

func (set *expiringSet) Remove(id string) {
	set.expire()
	set.values.Remove(id)
	delete(set.lastInserted, id)
}

func (set *expiringSet) Has(id string) (ok bool) {
	set.expire()
	ok = set.values.Has(id)
	return
}

func (set *expiringSet) Size() int {
	set.expire()
	return set.values.Size()
}

func (set *expiringSet) IDs() []string {
	set.expire()
	return set.values.IDs()
}

func (set *expiringSet) Values() map[string]struct{} {
	set.expire()
	return set.values.Values()
}

func (set *expiringSet) OnExpire(onExpire func(id string)) {
	set.onExpire = onExpire
}
