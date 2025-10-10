package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bsv-blockchain/libp2p-keygen/keygen"
)

const (
	defaultPort = 9095
	version     = "1.0.0"
)

func main() {
	// CLI flags
	vanityPattern := flag.String("vanity", "", "Search for a vanity peer ID containing this pattern")
	caseSensitive := flag.Bool("case-sensitive", false, "Enable case-sensitive pattern matching for vanity generation")
	verbose := flag.Bool("verbose", false, "Show verbose output during vanity generation")
	showVersion := flag.Bool("version", false, "Show version information")
	port := flag.Int("port", defaultPort, "Port for libp2p host display")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "libp2p-keygen v%s - Generate Ed25519 keys for libp2p\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  libp2p-keygen [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Generate a single key\n")
		fmt.Fprintf(os.Stderr, "  libp2p-keygen\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate a vanity key (case-insensitive)\n")
		fmt.Fprintf(os.Stderr, "  libp2p-keygen --vanity cool\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate a vanity key (case-sensitive)\n")
		fmt.Fprintf(os.Stderr, "  libp2p-keygen --vanity Cool --case-sensitive\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate with verbose output\n")
		fmt.Fprintf(os.Stderr, "  libp2p-keygen --vanity bitcoin --verbose\n\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("libp2p-keygen v%s\n", version)
		os.Exit(0)
	}

	var keyPair *keygen.KeyPair
	var err error

	// Vanity generation
	if *vanityPattern != "" {
		fmt.Printf("Searching for vanity peer ID containing '%s'", *vanityPattern)
		if *caseSensitive {
			fmt.Printf(" (case-sensitive)")
		}
		fmt.Printf("...\n")

		result, err := keygen.GenerateVanity(keygen.VanityOptions{
			Pattern:       *vanityPattern,
			CaseSensitive: *caseSensitive,
			Verbose:       *verbose,
		})
		if err != nil {
			log.Fatalf("Error generating vanity key: %v", err)
		}

		if *verbose {
			fmt.Println() // New line after progress output
		}

		keyPair = result.KeyPair

		fmt.Printf("\n✓ Found match after %d attempts (%.2f seconds, %.0f keys/sec)\n\n",
			result.Attempts,
			result.Duration.Seconds(),
			float64(result.Attempts)/result.Duration.Seconds(),
		)
	} else {
		// Simple generation
		keyPair, err = keygen.Generate()
		if err != nil {
			log.Fatalf("Error generating key: %v", err)
		}
	}

	// Output the key information
	fmt.Printf("╔═══════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                        libp2p Key Pair                            ║\n")
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════╝\n\n")

	fmt.Printf("Private Key (hex):\n%s\n\n", keyPair.PrivateKeyHex)
	fmt.Printf("Peer ID: %s\n\n", keyPair.PeerID.String())

	// Create a temporary host to display connection info
	libHost, err := keygen.CreateHost(keyPair.PrivateKey, *port)
	if err != nil {
		log.Fatalf("Error creating host: %v", err)
	}
	defer libHost.Close()

	fmt.Printf("Connect to this peer:\n")
	for _, addr := range libHost.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, libHost.ID().String())
	}

	// Teranode usage instructions
	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                    Using with Teranode                            ║\n")
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════╝\n\n")
	fmt.Printf("Add to your environment variables:\n")
	fmt.Printf("  export p2p_private_key=\"%s\"\n\n", keyPair.PrivateKeyHex)
	fmt.Printf("Or add to settings_local.conf:\n")
	fmt.Printf("  p2p_private_key = \"%s\"\n\n", keyPair.PrivateKeyHex)
	fmt.Printf("Note: The Peer ID is derived from your private key and does not need\n")
	fmt.Printf("      to be stored separately. Only the private key is required.\n")
}
