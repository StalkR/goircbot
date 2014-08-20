package battleroyale

import "testing"

func TestGet(t *testing.T) {
	t.Skip("leaderboard is down currently")
	s, err := get()
	if err != nil {
		t.Fatal(err)
	}
	for name, id := range map[string]string{
		"StalkR": "76561197960546028",
		"Ivan":   "76561197966750726",
	} {
		if _, ok := s[id]; !ok {
			t.Errorf("%s (%s): not found", name, id)
		}
	}
}
