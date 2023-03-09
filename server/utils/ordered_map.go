package utils

type OrderedMap[Key comparable, Value any] struct {
	innerMap map[Key]Value
	order    []Key
}

type KeyValuePair[Key comparable, Value any] struct {
	Key   Key
	Value Value
}

func (o OrderedMap[Key, Value]) Map() []KeyValuePair[Key, Value] {
	mapped := []KeyValuePair[Key, Value]{}

	for _, nextKey := range o.order {
		mapped = append(mapped, KeyValuePair[Key, Value]{
			Key:   nextKey,
			Value: o.innerMap[nextKey],
		})
	}

	return mapped
}

func (o OrderedMap[Key, Value]) Of(key Key) Value {
	return o.innerMap[key]
}

func (o OrderedMap[Key, Value]) HasKey(key Key) bool {
	_, hasKey := o.innerMap[key]
	return hasKey
}

func (o OrderedMap[Key, Value]) Upsert(key Key, value Value) OrderedMap[Key, Value] {
	if _, hasKey := o.innerMap[key]; !hasKey {
		o.order = append(o.order, key)
	}
	o.innerMap[key] = value

	return o
}

func NewOrderedMap[Key comparable, Value any]() OrderedMap[Key, Value] {
	return OrderedMap[Key, Value]{
		innerMap: map[Key]Value{},
		order:    []Key{},
	}
}
