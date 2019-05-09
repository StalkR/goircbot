package url

import (
	"testing"
)

func TestTwitter(t *testing.T) {
	for _, tt := range []struct {
		url  string
		want string
		err  error
	}{
		{url: "http://example.com", err: errSkip},
		{
			url:  "https://twitter.com/BenLaurie/status/331442973009133568",
			want: `Google Public DNS now checks DNSSEC for you by default. http://googleonlinesecurity.blogspot.co.uk/2013/03/google-public-dns-now-supports-dnssec.html ….`,
		},
		{
			url:  "https://twitter.com/supersat/status/331445098552369153",
			want: `@BenLaurie but no DNSSEC for http://google.com ? :(`,
		},
		{
			url:  "https://twitter.com/newsoft/status/484274141852622848",
			want: `@free_man_ Malheureusement @Ivanlef0u a eu une mauvaise influence sur @stalkr_ - maintenant il maitrise le #fapping comme personne :)`,
		},
		{
			url:  "https://twitter.com/element14/status/476395971472265216/photo/1",
			want: `There's a new #Arduino coming ! Have you seen it yet? http://ow.ly/xM7BJ https://pic.twitter.com/KV3hhRCi52`,
		},
		{
			url:  "https://twitter.com/DefConBeanBag1/status/761690424423571456/photo/1",
			want: "I want to break free... #DEFCON2016 @defcon @thedarktangent https://pic.twitter.com/egTrO0ja7q",
		},
	} {
		got, err := handleTwitter(tt.url)
		if tt.err != err {
			t.Errorf("Title(%v): err: %v", tt.url, err)
			continue
		}
		if err == nil && got != tt.want {
			t.Errorf("Title(%v): got %v; want %v", tt.url, got, tt.want)
		}
	}
}
