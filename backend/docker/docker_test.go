package docker

import "testing"

func TestGetImageName(t *testing.T) {
	got := GetImageName("test_path")
	want := "vsc-test_path-5da6ae5928d4a1ce395878ae9c7ea1f6"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}
