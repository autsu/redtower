package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(wd, "conf.yaml")
	cf, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", cf)
}

func TestInit(t *testing.T) {
	t.Logf("%+v\n", GlobalConf)
}
