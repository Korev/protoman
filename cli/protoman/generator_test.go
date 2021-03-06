package protoman

import (
	"os"
	"testing"
)

func TestGenerator(t *testing.T) {
	err := Generate("spotify.protoman", "registry", "test-directory/src/proto")
	if err != nil {
		t.Fatal("Failed to init protoman")
	}
	defer os.Remove(".protoman")
	defer os.RemoveAll("test-directory")
	if _, err := os.Stat(".protoman"); os.IsNotExist(err) {
		t.Error(".protoman does not exist")
	}
	if _, err := os.Stat("test-directory/src/proto/spotify/protoman/registry.proto"); os.IsNotExist(err) {
		t.Error("test-directory/src/proto/spotify/protoman/registry.proto does not exist")
	}
}
