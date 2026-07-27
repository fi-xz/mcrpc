package mcrpc

import (
	"crypto/tls"
	"time"
)

// An Option configures a Client. Options are applied by New and are
// order-independent: the TLS settings they touch are combined into a single
// configuration when Start dials.
type Option func(*Client)

// WithTLS makes the client connect over wss using cfg. A nil cfg selects TLS
// with default settings, which is what a server that presents a certificate but
// does not require client authentication needs.
//
// cfg is cloned, so later changes to your copy do not affect the client.
func WithTLS(cfg *tls.Config) Option {
	return func(c *Client) {
		c.useTLS = true
		c.tlsConfig = cfg
	}
}

// WithClientCertificate makes the client connect over wss and present cert for
// client certificate authentication. It implies WithTLS.
func WithClientCertificate(cert tls.Certificate) Option {
	return func(c *Client) {
		c.useTLS = true
		c.clientCert = &cert
	}
}

// WithInsecureSkipVerify disables verification of the server's certificate and
// implies WithTLS.
//
// This makes the connection vulnerable to interception: the client will accept
// any certificate the server offers, including one presented by an attacker.
// Use it only against servers with self-signed certificates on a network you
// control.
func WithInsecureSkipVerify() Option {
	return func(c *Client) {
		c.useTLS = true
		c.insecureSkipVerify = true
	}
}

// WithHandshakeTimeout caps how long Start waits for the WebSocket handshake.
// A non-positive duration removes the cap, leaving the context passed to Start
// as the only bound. Defaults to DefaultHandshakeTimeout.
func WithHandshakeTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.handshakeTimeout = d
	}
}

// WithTrace reports every JSON-RPC message the client sends and receives,
// with params and results exactly as they were serialised.
//
// This is a diagnostic aid for confirming what the client puts on the wire and
// what the server sends back — struct tags, omitted fields, and the JSON types
// of untyped values. It is called from the connection's read and write paths,
// so it must not block or call back into the client.
//
// Params may contain the server secret if a future method carries one; treat
// trace output as sensitive.
func WithTrace(trace func(TraceMessage)) Option {
	return func(c *Client) {
		c.trace = trace
	}
}

// WithHandler registers the callbacks invoked for server notifications. The
// handler is copied, so it must be complete at construction time; see Handler
// for why.
func WithHandler(h Handler) Option {
	return func(c *Client) {
		c.handler = h
	}
}

// tlsSettings combines the TLS-related options into one configuration, or
// returns nil for a plaintext connection.
func (c *Client) tlsSettings() *tls.Config {
	if !c.useTLS {
		return nil
	}

	cfg := &tls.Config{}
	if c.tlsConfig != nil {
		cfg = c.tlsConfig.Clone()
	}
	if c.insecureSkipVerify {
		cfg.InsecureSkipVerify = true
	}
	if c.clientCert != nil {
		cfg.Certificates = append(cfg.Certificates, *c.clientCert)
	}

	return cfg
}
