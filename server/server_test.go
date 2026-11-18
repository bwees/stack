package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"testing"

	"app/internal/controllers"

	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type sharedItem struct {
	Name string `json:"name"`
}

type sharedItemWrapper struct {
	Items []sharedItem `json:"items" description:"The wrapped items"`
}

type sharedItemList struct {
	Items []sharedItem `json:"items"`
}

type sharedItemsController struct{}

func (sharedItemsController) Register(group *fuego.Server) {
	fuego.Get(group, "/items", func(fuego.ContextNoBody) (sharedItemList, error) {
		return sharedItemList{}, nil
	})
}

type wrappedItemsController struct{}

func (wrappedItemsController) Register(group *fuego.Server) {
	fuego.Get(group, "/wrapped-items", func(fuego.ContextNoBody) (sharedItemWrapper, error) {
		return sharedItemWrapper{}, nil
	})
}

func newTestServer(t *testing.T, all ...controllers.Controller) *fuego.Server {
	t.Helper()
	return NewServer(fxtest.NewLifecycle(t), Controllers{In: fx.In{}, All: all})
}

func TestSpecIgnoresControllerOrder(t *testing.T) {
	spec := func(all ...controllers.Controller) string {
		t.Helper()
		raw, err := json.Marshal(newTestServer(t, all...).OutputOpenAPISpec())
		if err != nil {
			t.Fatalf("marshal spec: %v", err)
		}
		return string(raw)
	}

	if spec(sharedItemsController{}, wrappedItemsController{}) != spec(wrappedItemsController{}, sharedItemsController{}) {
		t.Fatal("spec depends on the order controllers register in")
	}
}

func TestServingDoesNotWriteSpec(t *testing.T) {
	t.Chdir(t.TempDir())

	newTestServer(t).OutputOpenAPISpec()

	if _, err := os.Stat("openapi.json"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("stat openapi.json = %v, want it not written", err)
	}
}
