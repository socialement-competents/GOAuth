package database

import "testing"

func TestConnecting(t *testing.T) {
	_, err := NewClient()
	if err != nil {
		t.Errorf("connecting the db failed: %v", err)
	}
}
