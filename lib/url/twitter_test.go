package url

import (
	"testing"
)

func TestTwitterMatch(t *testing.T) {
	p := &Twitter{}
	for _, tt := range []struct {
		url  string
		want bool
	}{
		{url: "http://twitter.com/nickname/status/123456", want: true},
		{url: "https://twitter.com/nickname/status/123456", want: true},
		{url: "http://twitter.com/nickname/status/123456#extra", want: true},
		{url: "http://example.com"},
	} {
		if got := p.Match(tt.url); got != tt.want {
			t.Errorf("Match(%s): got %v; want %v", tt.url, got, tt.want)
		}
	}
}

func TestTwitterParse(t *testing.T) {
	for _, tt := range []struct {
		url  string
		want string
	}{
		{
			url:  "https://twitter.com/BenLaurie/status/331442973009133568",
			want: `Google Public DNS now checks DNSSEC for you by default. http://googleonlinesecurity.blogspot.co.uk/2013/03/google-public-dns-now-supports-dnssec.html ….`,
		},
		{
			url:  "https://twitter.com/supersat/status/331445098552369153",
			want: `@BenLaurie but no DNSSEC for http://google.com ? :(`,
		},
		{
			url:  "https://twitter.com/newsoft/status/484274141852622848",
			want: `@free_man_ Malheureusement @Ivanlef0u a eu une mauvaise influence sur @stalkr_ - maintenant il maitrise le #fapping comme personne :)`,
		},
		{
			url:  "https://twitter.com/element14/status/476395971472265216/photo/1",
			want: `There's a new #Arduino coming ! Have you seen it yet? http://ow.ly/xM7BJ  pic.twitter.com/KV3hhRCi52`,
		},
	} {
		got, err := Title(tt.url)
		if err != nil {
			t.Errorf("Title(%s): err: %v", tt.url, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Title(%s): got %s; want %s", tt.url, got, tt.want)
		}
	}
}
