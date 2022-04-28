package docker

import (
	"os"
	"strings"
	"testing"
)

func TestSetVerboseTrue(t *testing.T) {
	d := new(Docker)
	d.SetVerbose(true)
	got := d.verbose
	want := true
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetVerboseFalse(t *testing.T) {
	d := new(Docker)
	d.SetVerbose(false)
	got := d.verbose
	want := false
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetPath(t *testing.T) {
	d := new(Docker)
	d.SetPath("test_path")
	got := d.path
	want := "test_path"
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetPathEmpty(t *testing.T) {
	d := new(Docker)
	d.SetPath("")
	got := d.path
	path, _ := os.Getwd()
	want := strings.ToLower(path)
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestGetPath(t *testing.T) {
	d := new(Docker)
	d.SetPath("test_path")
	got := d.GetPath()
	want := "test_path"
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetImageEmpty(t *testing.T) {
	d := new(Docker)
	d.SetPath("test_path")
	d.SetImage("")
	got := d.image
	want := "vsc-test_path-5da6ae5928d4a1ce395878ae9c7ea1f6"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestSetImageEmptyCapital(t *testing.T) {
	d := new(Docker)
	d.SetPath("Test_Path")
	d.SetImage("")
	got := d.image
	want := "vsc-test_path-5da6ae5928d4a1ce395878ae9c7ea1f6"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestSetImage(t *testing.T) {
	d := new(Docker)
	d.SetImage("test_image")
	got := d.image
	want := "test_image"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestGetImage(t *testing.T) {
	d := new(Docker)
	d.SetImage("test_image")
	got := d.GetImage()
	want := "test_image"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestSetContainer(t *testing.T) {
	d := new(Docker)
	d.SetContainer("test_container")
	got := d.container
	want := "test_container"
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetContainerEmpty(t *testing.T) {
	d := new(Docker)
	d.runOut = func([]string) (string, error) {
		return "test_container", nil
	}
	d.SetContainer("")
	got := d.container
	want := "test_container"
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestGetContainer(t *testing.T) {
	d := new(Docker)
	d.SetContainer("test_container")
	got := d.GetContainer()
	want := "test_container"
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetDockerfile(t *testing.T) {
	d := new(Docker)
	d.SetDockerfile("test_dockerfile")
	got := d.dockerfile
	want := "test_dockerfile"
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSetCommand(t *testing.T) {
	d := new(Docker)
	d.SetCommand([]string{"test", "command"})
	got := d.command
	want := []string{"test", "command"}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestSetArgs(t *testing.T) {
	d := new(Docker)
	d.SetArgs([]string{"test", "args"})
	got := d.args
	want := []string{"test", "args"}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestBuild(t *testing.T) {
	d := new(Docker)
	d.run = func(command []string, verbose bool) error {
		if len(command) < 2 {
			t.Errorf("missing command args")
		}
		return nil
	}
	if err := d.Build(); err != nil {
		t.Errorf("%v", err)
	}
}
