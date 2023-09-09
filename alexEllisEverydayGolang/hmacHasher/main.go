// Generate HMAC Digests for msg, using a key.
// Sample command to run:
//		go run main.go --msg="Hi." --key="1234"

package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		srcString string
		secretKey string
	)

	parseValidateInputs(&srcString, &secretKey)
	computeAndPrintHash(&srcString, &secretKey)
}

func computeAndPrintHash(srcString, secretKey *string) {
	fmt.Printf("Computing hash for Src string: %q\nSecret key: %q\n", *srcString, *secretKey)
	digest := signHmac([]byte(*srcString), []byte(*secretKey))
	fmt.Printf("Digest: %x\n", digest)
}

func parseValidateInputs(srcString, secretKey *string) {
	// Parse Inputs
	flag.StringVar(srcString, "msg", "", "Src string whose hash is required")
	flag.StringVar(secretKey, "key", "", "Non empty Secret key")
	flag.Parse()

	// Inout Validation
	if len(strings.TrimSpace(*srcString)) == 0 || len(strings.TrimSpace(*secretKey)) == 0 {
		fmt.Println("❌ Src string , or Key is empty or whitesspaced")
		os.Exit(1)
	}
}

func signHmac(msg, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	signed := mac.Sum(nil)
	return signed
}
