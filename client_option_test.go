package gax

import (
	"reflect"
	"testing"
)

func TestClientOptionsPieceByPiece(t *testing.T) {
	expected := &ClientSettings{
		"myapi",
		"v0.1.0",
		"https://example.com:443",
		[]string{"https://example.com/auth/helloworld", "https://example.com/auth/otherthing"},
	}

	settings := &ClientSettings{}
	opts := []ClientOption{
		WithAPIName("myapi"),
		WithAPIVersion("v0.1.0"),
		WithEndpoint("https://example.com:443"),
		WithScopes("https://example.com/auth/helloworld", "https://example.com/auth/otherthing"),
	}
	clientOptions(opts).Resolve(settings)

	if !reflect.DeepEqual(settings, expected) {
		t.Errorf("piece-by-piece settings don't match their expected configuration")
	}

	settings = &ClientSettings{}
	expected.Resolve(settings)

	if !reflect.DeepEqual(settings, expected) {
		t.Errorf("whole settings don't match their expected configuration")
	}

	expected.Scopes[0] = "hello"
	if settings.Scopes[0] == expected.Scopes[0] {
		t.Errorf("unexpected modification in Scopes array")
	}
}
