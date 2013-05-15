package url

import (
	"testing"
)

func TestGoodMatch(t *testing.T) {
	urls := []string{
		"http://twitter.com/nickname/status/123456",
		"https://twitter.com/nickname/status/123456",
		"http://twitter.com/nickname/status/123456#extra",
	}
	p := &Twitter{}
	for _, url := range urls {
		if !p.Match(url) {
			t.Errorf("Match(%v) = false, want true", url)
		}
	}
}

func TestBadMatch(t *testing.T) {
	urls := []string{
		"http://example.com",
	}
	p := &Twitter{}
	for _, url := range urls {
		if p.Match(url) {
			t.Errorf("Match(%v) = true, want false", url)
		}
	}
}

func TestParse1(t *testing.T) {
	const in = `      </div>

      <p class="js-tweet-text">Google Public DNS now checks DNSSEC for you by default. <a href="http://t.co/ZdWBummAXc" rel="nofollow" dir="ltr" data-expanded-url="http://googleonlinesecurity.blogspot.co.uk/2013/03/google-public-dns-now-supports-dnssec.html" class="twitter-timeline-link" target="_blank" title="http://googleonlinesecurity.blogspot.co.uk/2013/03/google-public-dns-now-supports-dnssec.html" ><span class="invisible">http://</span><span class="js-display-url">googleonlinesecurity.blogspot.co.uk/2013/03/google</span><span class="invisible">-public-dns-now-supports-dnssec.html</span><span class="tco-ellipsis">…</span></a>.</p>

      <div class="stream-item-footer">`
	const out = `Google Public DNS now checks DNSSEC for you by default. http://googleonlinesecurity.blogspot.co.uk/2013/03/google-public-dns-now-supports-dnssec.html….`
	p := &Twitter{}
	x, err := p.Parse(in)
	if err != nil {
		t.Error(err)
		return
	}
	if x != out {
		t.Errorf("Parse(%q) = %q (%v), want %q (%v)", in, x, len(x), out, len(out))
	}
}

func TestParse2(t *testing.T) {
	const in = `      </div>

      <p class="js-tweet-text tweet-text">Google app engine with PHP support ?! Does this mean you get bug bounties for PHP mem corruptions ? :D</p>

      <div class="stream-item-footer">`
	const out = `Google app engine with PHP support ?! Does this mean you get bug bounties for PHP mem corruptions ? :D`
	p := &Twitter{}
	x, err := p.Parse(in)
	if err != nil {
		t.Error(err)
		return
	}
	if x != out {
		t.Errorf("Parse(%q) = %q (%v), want %q (%v)", in, x, len(x), out, len(out))
	}
}
