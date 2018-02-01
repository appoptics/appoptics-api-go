package appoptics

import (
	"fmt"
	"net/url"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	token := "deadbeef"
	c := NewClient(token)

	t.Run("token should be set to passed-in value", func(t *testing.T) {
		if c.token != token {
			t.Errorf("expected '%s' to match '%s'", c.token, token)
		}
	})

	t.Run("baseURL should be set to default", func(t *testing.T) {
		clientURL, _ := url.Parse(defaultBaseURL)
		if *c.baseURL != *clientURL {
			t.Errorf("expected '%s' to match '%s'", *c.baseURL, *clientURL)
		}
	})

	t.Run("callerUserAgentFragment should be set to default", func(t *testing.T) {
		if c.completeUserAgentString() != clientVersionString() {
			t.Errorf("expected '%s' to match '%s'", c.completeUserAgentString(), clientVersionString())
		}
	})
}

func TestNewClient_Customized(t *testing.T) {
	token := "deadbeef"
	altUserAgentString := "totally-different-thing"
	altBaseURLString := "https://metrics-api.appoptics.com"

	t.Run("custom user agent string", func(t *testing.T) {
		c := NewClient(token, UserAgentClientOption(altUserAgentString))
		chkString := fmt.Sprintf("%s:%s", altUserAgentString, clientVersionString())
		req, _ := c.NewRequest("GET", "foo", nil)
		if req.UserAgent() != chkString {
			t.Errorf("expected '%s' to match '%s'", req.UserAgent(), chkString)
		}
	})

	t.Run("custom base URL", func(t *testing.T) {
		c := NewClient(token, BaseURLClientOption(altBaseURLString))
		testAltBaseURL, _ := url.Parse(altBaseURLString)
		if *c.baseURL != *testAltBaseURL {
			t.Errorf("expected '%s' to match '%s'", *c.baseURL, *testAltBaseURL)
		}
	})
}
