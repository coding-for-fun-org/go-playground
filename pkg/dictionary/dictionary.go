package dictionary

import "errors"

type Dictionary map[string]string

var (
	errNotFound   = errors.New("could not find the word you were looking for")
	errWordExists = errors.New("cannot add word because it already exists")
)

// Search method to search for a word in the dictionary
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]

	if exists {
		return value, nil
	}

	return "", errNotFound
}

// Add method to add a word to the dictionary
func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)
	if err == nil {
		return errWordExists
	}

	d[word] = definition

	return nil
}

// Update method to update a word in the dictionary
func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	if err != nil {
		return errNotFound
	}

	d[word] = definition

	return nil
}

// Delete method to delete a word from the dictionary
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
