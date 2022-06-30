// Package plex provides the methods the Notifiarr client uses to interface with Plex.
// This package also provides a web handler for incoming plex webhooks, and another
// two handlers for requests from Notifiarr.com to list sessions and kill a session.
// The purpose is to keep track of Plex viewers and send meaningful alerts to their
// respective Disord server about user behavior.
// ie. user started watching something, paused it, resumed it, and finished something.
// This package can be disabled by not providing a Plex Media Server URL or Token.
package plex

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Notifiarr/notifiarr/pkg/exp"
	"golift.io/cnfg"
)

// Server is the Plex configuration from a config file.
// Without a URL or Token, nothing works and this package is unused.
type Server struct {
	Timeout cnfg.Duration `toml:"timeout" json:"timeout" xml:"timeout"`
	URL     string        `toml:"url" json:"url" xml:"url"`
	Token   string        `toml:"token" json:"token" xml:"token"`
	Name    string        `toml:"-" json:"-" xml:"-"`
	client  *http.Client
}

const (
	defaultTimeout = 10 * time.Second
	minimumTimeout = 2 * time.Second
)

// ErrNoURLToken is returned when there is no token or URL.
var ErrNoURLToken = fmt.Errorf("token or URL for Plex missing")

// Configured returns true ifthe server is configured, false otherwise.
func (s *Server) Configured() bool {
	return s != nil && s.URL != "" && s.Token != ""
}

// Validate checks input values and starts the cron interval if it's configured.
func (s *Server) Validate() {
	s.URL = strings.TrimRight(s.URL, "/")

	if s.Timeout.Duration == 0 {
		s.Timeout.Duration = defaultTimeout
	} else if s.Timeout.Duration < minimumTimeout {
		s.Timeout.Duration = minimumTimeout
	}
}

func (s *Server) getPlexURL(ctx context.Context, url string, params url.Values) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("X-Plex-Token", s.Token)
	req.Header.Set("Accept", "application/json")
	exp.Apps.Add("Plex&&GET Requests", 1)

	resp, err := s.getClient().Do(req)
	if err != nil {
		exp.Apps.Add("Plex&&GET Errors", 1)
		return nil, fmt.Errorf("making http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exp.Apps.Add("Plex&&GET Errors", 1)
		return nil, fmt.Errorf("reading http response: %w", err)
	}

	exp.Apps.Add("Plex&&Bytes Received", int64(len(body)))

	if resp.StatusCode != http.StatusOK {
		exp.Apps.Add("Plex&&GET Errors", 1)
		return body, ErrBadStatus
	}

	return body, nil
}

func (s *Server) getClient() *http.Client {
	if s.client == nil {
		s.client = &http.Client{}
	}

	return s.client
}
