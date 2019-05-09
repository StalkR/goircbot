package url

import (
  "testing"
)

func TestDefault(t *testing.T) {
  for _, tt := range []struct {
    url  string
    want string
  }{
    {
      url:  "https://stalkr.net/",
      want: "stalkr.net",
    },
  } {
    got, err := handleDefault(tt.url)
    if err != nil {
      t.Errorf("Title(%v): err: %v", tt.url, err)
      continue
    }
    if got != tt.want {
      t.Errorf("Title(%v): got %v; want %v", tt.url, got, tt.want)
    }
  }
}
