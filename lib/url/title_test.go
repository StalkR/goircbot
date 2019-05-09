package url

import (
  "testing"
)

func TestTitle(t *testing.T) {
  badCalled := false
  origHandlers := handlers
  handlers = []func(url string) (string, error){
    func(url string) (string, error) {
      badCalled = true
      return "bad", errSkip
    },
    func(url string) (string, error) {
      return "good", nil
    },
  }
  defer func() {
    handlers = origHandlers
  }()

  got, err := Title("https://example.com/")
  if err != nil {
    t.Fatal(err)
  }
  if want := "good"; got != want {
    t.Errorf("got %v, want %v", got, want)
  }
  if !badCalled {
    t.Error("expected bad to be called but did not")
  }
}
