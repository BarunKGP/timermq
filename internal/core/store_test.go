package core

import "testing"

func TestStore(t *testing.T) {
	store := NewStore[string]()

	for _, msg := range []string{"terrier", "dalmatian", "retriever"} {
		_ = store.Consume(msg)
	}

	if store.Len() != 3 {
		t.Errorf("Unexpected items in store. Expected 3, found %d", store.Len())
	}

	msg2, err := store.Get(2)
	if err != nil {
		t.Error(err)
	}

	if msg2 != "retriever" {
		t.Errorf("Unexpected message retrieved. Expected \"retriever\", found %s", msg2)
	}
}
