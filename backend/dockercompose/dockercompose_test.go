package dockercompose

import (
	"testing"
)

func TestSetCommand(t *testing.T) {
	got := setCommand("project", "file", "action")
	want := []string{"docker-compose", "--project-name", "project", "--file", "file", "action"}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
}

func TestSetCommandArgs(t *testing.T) {
	got := setCommand("project", "file", "action", "arg1", "arg2")
	want := []string{"docker-compose", "--project-name", "project", "--file", "file", "action", "arg1", "arg2"}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
}
