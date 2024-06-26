package common

import "golang.org/x/exp/constraints"

func ListContainApi[T constraints.Ordered](list *[]T, entity T) bool {
	if list == nil || len(*list) == 0 {
		return false
	}
	for _, element := range *list {
		if element == entity {
			return true
		}
	}
	return false
}

func GetMapKey[K constraints.Ordered, V any](data *map[K]V) *[]K {
	resp := make([]K, 0)
	if data == nil || len(*data) == 0 {
		return &resp
	}
	for key, _ := range *data {
		resp = append(resp, key)
	}
	return &resp
}

func GetMapValue[K constraints.Ordered, V any](data *map[K]V) *[]V {
	resp := make([]V, 0)
	if data == nil || len(*data) == 0 {
		return &resp
	}
	for _, value := range *data {
		resp = append(resp, value)
	}
	return &resp
}

func GetMapValueDistinct[K constraints.Ordered, V constraints.Ordered](data *map[K]V) *[]V {
	resp := make([]V, 0)
	if data == nil || len(*data) == 0 {
		return &resp
	}
	m := make(map[V]string, 0)
	for _, value := range *data {
		m[value] = ""
	}
	for value, _ := range m {
		resp = append(resp, value)
	}
	return &resp
}

func ListDistinct[T constraints.Ordered](data *[]T) *[]T {
	resp := make([]T, 0)
	if data == nil || len(*data) == 0 {
		return &resp
	}
	m := make(map[T]string)
	for _, t := range *data {
		m[t] = ""
	}
	for t, _ := range m {
		resp = append(resp, t)
	}
	return &resp
}
