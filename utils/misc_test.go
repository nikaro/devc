package utils

import "testing"

func TestRemoveFromSlice(t *testing.T) {
	got := RemoveFromSlice([]string{"a", "b", "c", "c"}, "c")
	want := []string{"a", "b"}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
}

func TestMd5Sum(t *testing.T) {
	got := Md5Sum("testing")
	want := "ae2b1fca515949e5d54fb22b8ed95575"
	if got != want {
		t.Errorf("got: %s. want: %s", got, want)
	}
}
