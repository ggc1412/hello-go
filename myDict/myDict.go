package myDict

import "errors"

type Dictionary map[string]string

var (
	errNotFound   = errors.New("Not Found")
	errCantUpdate = errors.New("Cant update non-existing word")
	errWordExists = errors.New("That word already exists")
)

func (d Dictionary) Search(word string) (string, error) {
	word, exist := d[word]
	if exist {
		return word, nil
	}
	return "", errNotFound
}

func (d Dictionary) Add(newWord, definition string) error {
	_, err := d.Search(newWord)
	switch err {
	case errNotFound:
		d[newWord] = definition
	case nil:
		return errWordExists
	}
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case errNotFound:
		return errNotFound
	}
	return nil
}
