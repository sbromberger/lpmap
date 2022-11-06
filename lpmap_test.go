package lpmap

import (
	"fmt"
	"math/rand"
	"testing"
)

type myKey uint32

func (k myKey) Hash() uint64 {
	return uint64(k)
}

type myCollKey uint32

func (k myCollKey) Hash() uint64 {
	return uint64(k) % 5
}

var mapSizes = []int{10_000, 100_000, 1_000_000, 10_000_000, 100_000_000}
var thresholds = []float64{0.5}

var size = 10

func makeMap(threshold float64) (Map[myKey, uint64], []myKey) {
	keys := make([]myKey, size)
	for i := 0; i < size; i++ {
		keys[i] = myKey(rand.Uint64() % 100000)
	}
	lp := New[myKey, uint64](size, threshold)

	for i := 0; i < size; i++ {
		v := uint64(i)
		lp.Set(keys[i], v)
	}
	return lp, keys
}

func TestGet(t *testing.T) {
	for _, threshold := range thresholds {
		lp, keys := makeMap(threshold)

		t.Run("Get valid keys", func(t *testing.T) {
			for i, k := range keys {
				v, found := lp.Get(k)
				if !found {
					t.Errorf("Key %v should be found but wasn't", k)
				}
				if *v != uint64(i) {
					t.Errorf("Value %d for key %v does not match expected value %d", *v, k, i)

				}
			}
		})
		t.Run("Get invalid keys", func(t *testing.T) {
			for i := 100000; i < 1_000_000; i += 100000 {
				k := myKey(i)
				v, found := lp.Get(k)
				if found {
					t.Errorf("Key %v should not be found but was, with value %d", k, v)
				}
			}
		})
	}
}

func TestSet(t *testing.T) {
	size := 10
	for _, threshold := range thresholds {
		lp := New[myCollKey, uint64](0, threshold)
		for i := 0; i < size; i++ {
			lp.Set(myCollKey(i), uint64(i+size))
		}

		if lp.Size() != size {
			t.Errorf("numEntries should be %d; got %d\n", size, lp.Size())
		}

		for i := 0; i < size+1; i++ {
			lp.Set(myCollKey(i), uint64(i))
		}
		if lp.Size() != size+1 {
			t.Errorf("numEntries should be %d; got %d (%+v)\n", size+1, lp.Size(), lp)
		}

	}
}

func TestDelete(t *testing.T) {
	size := 10
	for _, threshold := range thresholds {
		lp := New[myCollKey, uint64](0, threshold)
		for i := 0; i < size; i++ {
			lp.Set(myCollKey(i), uint64(i+size))
		}

		for i := 0; i < size; i += 2 {
			found := lp.Delete(myCollKey(i))
			if !found {
				t.Errorf("Tried to delete key %v but it wasn't found", myCollKey(i))
			}
		}
		if lp.Size() != size/2 {
			t.Errorf("Post-delete numEntries should be %d; got %d\n", size/2, lp.Size())
		}

		for i := 0; i < size; i += 2 {
			found := lp.Delete(myCollKey(i))
			if found {
				t.Errorf("Deleted nonexistent key %v, lp is %+v", myCollKey(i), lp)
			}
		}
		if lp.Size() != size/2 {
			t.Errorf("Post-delete numEntries should be %d; got %d\n", size/2, lp.Size())
		}
	}

}
func BenchmarkGet(b *testing.B) {
	for _, threshold := range thresholds {
		for _, size := range mapSizes {
			lp := New[myKey, uint64](0, threshold)
			m := make(map[myKey]uint64)
			for i := 0; i < size; i++ {
				r := rand.Uint64()
				k := myKey(r)
				v := uint64(i)
				lp.Set(k, v)
				m[k] = v
			}

			b.ResetTimer()

			b.Run(fmt.Sprintf("Get/lp/%d @ %f", size, threshold), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					lp.Get(myKey(i))
				}
			})
			b.Run(fmt.Sprintf("Get/map/%d", size), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = m[myKey(i)]
				}
			})
		}

	}
}

func BenchmarkSet(b *testing.B) {
	for _, threshold := range thresholds {
		b.ResetTimer()
		for _, size := range mapSizes {
			lp := New[myKey, uint64](0, threshold)
			m := make(map[myKey]uint64)

			b.ResetTimer()

			b.Run(fmt.Sprintf("Set/lp/%d @ %f", size, threshold), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					lp.Set(myKey(rand.Uint64()), uint64(i))
				}
			})

			b.Run(fmt.Sprintf("Set/map/%d @ %f", size, threshold), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					m[myKey(rand.Uint64())] = uint64(i)
				}
			})

		}

	}
}
