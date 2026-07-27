package mcrpc

import (
	"crypto/tls"
	"testing"
	"time"
)

func TestPlaintextByDefault(t *testing.T) {
	if config := New("host:1", "secret").tlsSettings(); config != nil {
		t.Errorf("a client with no TLS option should dial in plaintext, got %+v", config)
	}
}

func TestTLSOptionsCombine(t *testing.T) {
	certificate := tls.Certificate{Certificate: [][]byte{{0x01}}}

	// The options record intent and Start assembles the configuration, so the
	// order they are given in must not matter.
	orders := map[string][]Option{
		"tls first": {
			WithTLS(&tls.Config{ServerName: "example.test"}),
			WithClientCertificate(certificate),
			WithInsecureSkipVerify(),
		},
		"tls last": {
			WithInsecureSkipVerify(),
			WithClientCertificate(certificate),
			WithTLS(&tls.Config{ServerName: "example.test"}),
		},
	}

	for name, options := range orders {
		t.Run(name, func(t *testing.T) {
			config := New("host:1", "secret", options...).tlsSettings()
			if config == nil {
				t.Fatal("expected a TLS configuration")
			}
			if config.ServerName != "example.test" {
				t.Errorf("ServerName = %q, want example.test", config.ServerName)
			}
			if !config.InsecureSkipVerify {
				t.Error("InsecureSkipVerify was not applied")
			}
			if len(config.Certificates) != 1 {
				t.Errorf("got %d certificates, want 1", len(config.Certificates))
			}
		})
	}
}

func TestWithTLSClonesTheConfig(t *testing.T) {
	// The caller keeps their own config, and neither side may reach the other
	// through it once the option has been applied.
	t.Run("a later change by the caller does not reach the client", func(t *testing.T) {
		supplied := &tls.Config{ServerName: "example.test"}
		client := New("host:1", "secret", WithTLS(supplied))

		supplied.ServerName = "mutated"
		supplied.InsecureSkipVerify = true

		config := client.tlsSettings()
		if config.ServerName != "example.test" {
			t.Errorf("ServerName = %q, want the value given to WithTLS", config.ServerName)
		}
		if config.InsecureSkipVerify {
			t.Error("the caller turned verification off after the fact")
		}
	})

	t.Run("the assembled config is not the caller's", func(t *testing.T) {
		supplied := &tls.Config{ServerName: "example.test"}
		client := New("host:1", "secret", WithTLS(supplied))

		config := client.tlsSettings()
		config.ServerName = "mutated"

		if supplied.ServerName != "example.test" {
			t.Errorf("the caller's config was mutated: ServerName = %q", supplied.ServerName)
		}
	})
}

func TestWithTLSNilSelectsDefaults(t *testing.T) {
	// A server that presents a certificate but does not ask for one needs TLS
	// with nothing configured, which must still select wss.
	client := New("host:1", "secret", WithTLS(nil))

	if !client.useTLS {
		t.Error("WithTLS(nil) should still select TLS")
	}
	if config := client.tlsSettings(); config == nil {
		t.Error("expected a TLS configuration")
	}
}

func TestWithHandshakeTimeout(t *testing.T) {
	if got := New("host:1", "secret").handshakeTimeout; got != DefaultHandshakeTimeout {
		t.Errorf("default handshake timeout = %s, want %s", got, DefaultHandshakeTimeout)
	}

	client := New("host:1", "secret", WithHandshakeTimeout(3*time.Second))
	if client.handshakeTimeout != 3*time.Second {
		t.Errorf("handshake timeout = %s, want 3s", client.handshakeTimeout)
	}

	// A non-positive duration removes the cap, leaving the context as the only
	// bound.
	client = New("host:1", "secret", WithHandshakeTimeout(0))
	if client.handshakeTimeout != 0 {
		t.Errorf("handshake timeout = %s, want 0", client.handshakeTimeout)
	}
}
