package models

import (
	"encoding/json"
	"testing"
)

func TestMarshallingUser(t *testing.T) {
	u := User{
		ID:       123456,
		Provider: "test",
		GHUser: &GHUser{
			Name:  "miguel",
			Login: "potato",
		},
	}
	ref := &u

	bytes, err := json.Marshal(ref)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(bytes))
	}

	var v User
	if err = json.Unmarshal(bytes, &v); err != nil {
		t.Error(err)
	}

	if u.ID != v.ID || u.Provider != v.Provider {
		t.Error("u and v don't match")
	}

	if u.Name != v.GHUser.Name || u.GHUser.Login != v.Login {
		t.Error("it doesn't marshal ghuser")
	}

	if u.ID == 0 || u.Provider == "" || u.Login == "" {
		t.Error("info loss")
	}
}
