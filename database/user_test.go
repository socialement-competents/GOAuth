package database

import (
	"testing"
)

func TestGettingInexistingUser(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("connecting the db failed: %v", err)
	}

	u, err := client.GetUserByLogin("does not exist", "still does not exist")
	t.Log(u, err)
}
