package collection

import "sync"

type Value interface{}

type ThreadSafeArray struct {
	items []Value
	lock  sync.Mutex
}

func (array *ThreadSafeArray) StartTransaction() {
	array.lock.Lock()
}

func (array *ThreadSafeArray) EndTransaction() {
	array.lock.Unlock()
}

func (array *ThreadSafeArray) Add(value Value) {
	if array.items == nil {
		array.items = []Value{}
	}
	array.items = append(array.items, value)
}

func (array *ThreadSafeArray) Remove(index int) {
	if index < len(array.items) {
		array.items = append(array.items[:index], array.items[index+1:]...)
	}
}

func (array *ThreadSafeArray) Get(index int) Value {
	if index < len(array.items) {
		return array.items[index]
	}
	return nil
}

func (array *ThreadSafeArray) Size() int {
	return len(array.items)
}
