package keygen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// KeyPair represents a generated key pair with its peer ID
type KeyPair struct {
	PrivateKeyHex string
	PublicKeyHex  string
	PrivateKey    crypto.PrivKey
	PublicKey     crypto.PubKey
	PeerID        peer.ID
}

// Generate creates a single Ed25519 key pair
func Generate() (*KeyPair, error) {
	privateKey, publicKey, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	privateRaw, err := privateKey.Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw private key: %w", err)
	}

	publicRaw, err := publicKey.Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw public key: %w", err)
	}

	privateHex := hex.EncodeToString(privateRaw)
	publicHex := hex.EncodeToString(publicRaw)

	peerID, err := peer.IDFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive peer ID: %w", err)
	}

	return &KeyPair{
		PrivateKeyHex: privateHex,
		PublicKeyHex:  publicHex,
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		PeerID:        peerID,
	}, nil
}

// VanityOptions configures vanity peer ID generation
type VanityOptions struct {
	Pattern       string
	CaseSensitive bool
	Verbose       bool
}

// VanityResult contains the result of a vanity search
type VanityResult struct {
	KeyPair  *KeyPair
	Attempts int
	Duration time.Duration
}

// GenerateVanity continuously generates keys until finding one matching the pattern
func GenerateVanity(opts VanityOptions) (*VanityResult, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("vanity pattern cannot be empty")
	}

	pattern := opts.Pattern
	if !opts.CaseSensitive {
		pattern = strings.ToLower(pattern)
	}

	startTime := time.Now()
	attempts := 0

	for {
		attempts++

		keyPair, err := Generate()
		if err != nil {
			return nil, err
		}

		peerIDStr := keyPair.PeerID.String()
		searchStr := peerIDStr
		if !opts.CaseSensitive {
			searchStr = strings.ToLower(peerIDStr)
		}

		if strings.Contains(searchStr, pattern) {
			return &VanityResult{
				KeyPair:  keyPair,
				Attempts: attempts,
				Duration: time.Since(startTime),
			}, nil
		}

		if opts.Verbose && attempts%1000 == 0 {
			fmt.Printf("\rAttempts: %d (%.2f keys/sec)", attempts, float64(attempts)/time.Since(startTime).Seconds())
		}
	}
}

// DecodePrivateKey decodes a hex-encoded Ed25519 private key
func DecodePrivateKey(hexKey string) (crypto.PrivKey, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex key: %w", err)
	}

	privateKey, err := crypto.UnmarshalEd25519PrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %w", err)
	}

	return privateKey, nil
}

// CreateHost creates a libp2p host from a private key for testing/display
func CreateHost(privateKey crypto.PrivKey, port int) (host.Host, error) {
	libHost, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	return libHost, nil
}
