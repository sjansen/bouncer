package config

import "net/url"

// URL represents a parsed URL
type URL struct {
	url.URL
}

// UnmarshalText converts an environment variable string to a URL
func (u *URL) UnmarshalText(text []byte) error {
	return u.URL.UnmarshalBinary(text)
}

// ResolveReference proxies to net/url.URL.ResolveReference()
func (u *URL) ResolveReference(ref *url.URL) *URL {
	return &URL{
		URL: *u.URL.ResolveReference(ref),
	}
}
