package url

// Parsers contains an ordered list of supported URL title Parsers.
var Parsers = []Parser{
	&Twitter{},
	&Default{},
}

// A Parser can Parse() a body to extract a title if a given URL Match()es.
type Parser interface {
	Match(url string) bool
	Parse(body string) (string, error)
}

// ParseTitle parses a title from an URL and content given a list of Parsers.
func ParseTitle(url, body string, parsers []Parser) (string, error) {
	for _, p := range parsers {
		if p.Match(url) {
			return p.Parse(body)
		}
	}
	panic("no parser matched")
}
