package kel

// Object ...
type Object interface {
	GetResourceType() string
	GetID() string
}

// Store represents objects retrieved from the Kel API
type Store struct {
	objs       map[string][]Object
	sortedObjs map[string]map[string][]Object
}

// AddSortedList ..
func (store *Store) AddSortedList(key string, objs []Object) {
}

// Add ...
func (store *Store) Add(obj Object) {
}

// ListSorted ...
func (store *Store) ListSorted(resourceType string, key string) []Object {
	return nil
}

// List ...
func (store *Store) List(resourceType string) []Object {
	return nil
}

// Get ...
func (store *Store) Get(resourceType string, ID string) Object {
	return nil
}
