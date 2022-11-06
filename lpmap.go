package lpmap

type KeyType interface {
	Hash() uint64
	comparable
}

type status uint8

const (
	dead     status = 1
	occupied status = 2
)

type Map[K KeyType, V any] struct {
	keys       []K
	values     []V
	statuses   []status
	threshold  float64
	numEntries int
}

func New[K KeyType, V any](size int, threshold float64) Map[K, V] {
	nEntries := int(float64(size)/threshold) + 1
	keys := make([]K, nEntries)
	values := make([]V, nEntries)
	statuses := make([]status, nEntries)
	return Map[K, V]{keys, values, statuses, threshold, 0}
}

func (m *Map[K, V]) Get(k K) (*V, bool) {
	if m.numEntries == 0 {
		return nil, false
	}
	i := int(k.Hash() % uint64(len(m.keys)))
	var coll int
	for {
		status := m.statuses[i]
		if status == occupied && m.keys[i] == k {
			return &m.values[i], true
		}
		if status == 0 {
			return nil, false
		}
		i++
		coll++
		if i == len(m.keys) {
			i = 0
		}
	}
}

func (m *Map[K, V]) Values() chan V {
	ch := make(chan V, m.numEntries)
	defer close(ch)
	go func() {
		for i, v := range m.values {
			if m.statuses[i] == occupied {
				ch <- v
			}
		}
	}()

	return ch

}
func getNextAvailableIndex[K KeyType](keys []K, statuses []status, k K) int {
	i := int(k.Hash() % uint64(len(keys)))

	for {
		status := statuses[i]
		if status != occupied {
			return i
		}
		i++
		if i == len(keys) {
			i = 0
		}
	}
}
func (m *Map[K, V]) resize(newSize int) {

	if newSize < m.numEntries+1 {
		newSize = m.numEntries + 1
	}

	newKeys := make([]K, newSize)
	newValues := make([]V, newSize)
	newStatuses := make([]status, newSize)
	var count int
	for i, k := range m.keys {
		if m.statuses[i] == occupied {
			newI := getNextAvailableIndex(newKeys, newStatuses, k)
			newKeys[newI] = k
			newValues[newI] = m.values[i]
			newStatuses[newI] = occupied
			count++
		}
	}
	newMap := Map[K, V]{
		keys:       newKeys,
		values:     newValues,
		statuses:   newStatuses,
		numEntries: count,
		threshold:  m.threshold,
	}
	*m = newMap
}

func (m *Map[K, V]) Set(k K, v V) {
	if float64(m.numEntries)+1 > float64(len(m.keys))*m.threshold {
		m.resize(2 * len(m.keys))
	}
	i := k.Hash() % uint64(len(m.keys))
	for {
		status := m.statuses[i]
		if status != occupied {
			m.keys[i] = k
			m.values[i] = v
			m.statuses[i] = occupied
			m.numEntries++
			return
		}
		if status == occupied {
			if m.keys[i] == k {
				m.values[i] = v
				return
			}
		}
		i++
		if i == uint64(len(m.keys)) {
			i = 0
		}
	}
}

func (m *Map[K, V]) Delete(k K) bool {
	if m.numEntries == 0 {
		return false
	}
	i := int(k.Hash() % uint64(len(m.keys)))
	for {
		status := m.statuses[i]
		if status == occupied && m.keys[i] == k {
			m.statuses[i] = dead
			m.numEntries--
			return true
		}
		if status == 0 {
			return false
		}
		i++
		if i == len(m.keys) {
			i = 0
		}
	}
}

func (m *Map[K, V]) Size() int {
	return m.numEntries
}
