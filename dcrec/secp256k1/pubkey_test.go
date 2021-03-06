// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package secp256k1

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type pubKeyTest struct {
	name    string
	key     []byte
	format  byte
	isValid bool
}

var pubKeyTests = []pubKeyTest{
	// pubkey from bitcoin blockchain tx
	// 0437cd7f8525ceed2324359c2d0ba26006d92d85
	{
		name: "uncompressed ok",
		key: []byte{0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: true,
		format:  pubkeyUncompressed,
	},
	{
		name: "uncompressed x changed",
		key: []byte{0x04, 0x15, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: false,
	},
	{
		name: "uncompressed y changed",
		key: []byte{0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa4,
		},
		isValid: false,
	},
	{
		name: "uncompressed claims compressed",
		key: []byte{0x03, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: false,
	},
	{
		name: "uncompressed as hybrid ok",
		key: []byte{0x07, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: false,
	},
	{
		name: "uncompressed as hybrid wrong",
		key: []byte{0x06, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: false,
	},
	// from tx 0b09c51c51ff762f00fb26217269d2a18e77a4fa87d69b3c363ab4df16543f20
	{
		name: "compressed ok (ybit = 0)",
		key: []byte{0x02, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
			0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75, 0xf1,
			0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b, 0xd0, 0x21,
			0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b, 0x4d,
		},
		isValid: true,
		format:  pubkeyCompressed,
	},
	// from tx fdeb8e72524e8dab0da507ddbaf5f88fe4a933eb10a66bc4745bb0aa11ea393c
	{
		name: "compressed ok (ybit = 1)",
		key: []byte{0x03, 0x26, 0x89, 0xc7, 0xc2, 0xda, 0xb1, 0x33,
			0x09, 0xfb, 0x14, 0x3e, 0x0e, 0x8f, 0xe3, 0x96, 0x34,
			0x25, 0x21, 0x88, 0x7e, 0x97, 0x66, 0x90, 0xb6, 0xb4,
			0x7f, 0x5b, 0x2a, 0x4b, 0x7d, 0x44, 0x8e,
		},
		isValid: true,
		format:  pubkeyCompressed,
	},
	{
		name: "compressed claims uncompressed (ybit = 0)",
		key: []byte{0x04, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
			0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75, 0xf1,
			0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b, 0xd0, 0x21,
			0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b, 0x4d,
		},
		isValid: false,
	},
	{
		name: "compressed claims uncompressed (ybit = 1)",
		key: []byte{0x05, 0x26, 0x89, 0xc7, 0xc2, 0xda, 0xb1, 0x33,
			0x09, 0xfb, 0x14, 0x3e, 0x0e, 0x8f, 0xe3, 0x96, 0x34,
			0x25, 0x21, 0x88, 0x7e, 0x97, 0x66, 0x90, 0xb6, 0xb4,
			0x7f, 0x5b, 0x2a, 0x4b, 0x7d, 0x44, 0x8e,
		},
		isValid: false,
	},
	{
		name:    "wrong length)",
		key:     []byte{0x05},
		isValid: false,
	},
	{
		name: "X == P",
		key: []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFC, 0x2F, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: false,
	},
	{
		name: "X > P",
		key: []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFD, 0x2F, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: false,
	},
	{
		name: "Y == P",
		key: []byte{0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF,
			0xFF, 0xFC, 0x2F,
		},
		isValid: false,
	},
	{
		name: "Y > P",
		key: []byte{0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF,
			0xFF, 0xFD, 0x2F,
		},
		isValid: false,
	},
	{
		name: "hybrid",
		key: []byte{0x06, 0x79, 0xbe, 0x66, 0x7e, 0xf9, 0xdc, 0xbb,
			0xac, 0x55, 0xa0, 0x62, 0x95, 0xce, 0x87, 0x0b, 0x07,
			0x02, 0x9b, 0xfc, 0xdb, 0x2d, 0xce, 0x28, 0xd9, 0x59,
			0xf2, 0x81, 0x5b, 0x16, 0xf8, 0x17, 0x98, 0x48, 0x3a,
			0xda, 0x77, 0x26, 0xa3, 0xc4, 0x65, 0x5d, 0xa4, 0xfb,
			0xfc, 0x0e, 0x11, 0x08, 0xa8, 0xfd, 0x17, 0xb4, 0x48,
			0xa6, 0x85, 0x54, 0x19, 0x9c, 0x47, 0xd0, 0x8f, 0xfb,
			0x10, 0xd4, 0xb8,
		},
		isValid: false,
	},
}

func TestPubKeys(t *testing.T) {
	for _, test := range pubKeyTests {
		pk, err := ParsePubKey(test.key)
		if err != nil {
			if test.isValid {
				t.Errorf("%s pubkey failed when shouldn't %v",
					test.name, err)
			}
			continue
		}
		if !test.isValid {
			t.Errorf("%s counted as valid when it should fail",
				test.name)
			continue
		}
		var pkStr []byte
		switch test.format {
		case pubkeyUncompressed:
			pkStr = pk.SerializeUncompressed()
		case pubkeyCompressed:
			pkStr = pk.SerializeCompressed()
		}
		if !bytes.Equal(test.key, pkStr) {
			t.Errorf("%s pubkey: serialized keys do not match.",
				test.name)
			spew.Dump(test.key)
			spew.Dump(pkStr)
		}
	}
}

func TestPublicKeyIsEqual(t *testing.T) {
	pubKey1, err := ParsePubKey(
		[]byte{0x03, 0x26, 0x89, 0xc7, 0xc2, 0xda, 0xb1, 0x33,
			0x09, 0xfb, 0x14, 0x3e, 0x0e, 0x8f, 0xe3, 0x96, 0x34,
			0x25, 0x21, 0x88, 0x7e, 0x97, 0x66, 0x90, 0xb6, 0xb4,
			0x7f, 0x5b, 0x2a, 0x4b, 0x7d, 0x44, 0x8e,
		},
	)
	if err != nil {
		t.Fatalf("failed to parse raw bytes for pubKey1: %v", err)
	}

	pubKey2, err := ParsePubKey(
		[]byte{0x02, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
			0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75, 0xf1,
			0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b, 0xd0, 0x21,
			0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b, 0x4d,
		},
	)
	if err != nil {
		t.Fatalf("failed to parse raw bytes for pubKey2: %v", err)
	}

	if !pubKey1.IsEqual(pubKey1) {
		t.Fatalf("value of IsEqual is incorrect, %v is "+
			"equal to %v", pubKey1, pubKey1)
	}

	if pubKey1.IsEqual(pubKey2) {
		t.Fatalf("value of IsEqual is incorrect, %v is not "+
			"equal to %v", pubKey1, pubKey2)
	}
}

// TestDecompressY ensures that decompressY works as expected for some edge
// cases.
func TestDecompressY(t *testing.T) {
	tests := []struct {
		name      string // test description
		x         string // hex encoded x coordinate
		valid     bool   // expected decompress result
		wantOddY  string // hex encoded expected odd y coordinate
		wantEvenY string // hex encoded expected even y coordinate
	}{{
		name:      "x = 0 -- not a point on the curve",
		x:         "0",
		valid:     false,
		wantOddY:  "",
		wantEvenY: "",
	}, {
		name:      "x = 1",
		x:         "1",
		valid:     true,
		wantOddY:  "bde70df51939b94c9c24979fa7dd04ebd9b3572da7802290438af2a681895441",
		wantEvenY: "4218f20ae6c646b363db68605822fb14264ca8d2587fdd6fbc750d587e76a7ee",
	}, {
		name:      "x = secp256k1 prime (aka 0) -- not a point on the curve",
		x:         "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f",
		valid:     false,
		wantOddY:  "",
		wantEvenY: "",
	}, {
		name:      "x = secp256k1 prime - 1 -- not a point on the curve",
		x:         "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2e",
		valid:     false,
		wantOddY:  "",
		wantEvenY: "",
	}, {
		name:      "x = secp256k1 group order",
		x:         "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141",
		valid:     true,
		wantOddY:  "670999be34f51e8894b9c14211c28801d9a70fde24b71d3753854b35d07c9a11",
		wantEvenY: "98f66641cb0ae1776b463ebdee3d77fe2658f021db48e2c8ac7ab4c92f83621e",
	}, {
		name:      "x = secp256k1 group order - 1 -- not a point on the curve",
		x:         "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140",
		valid:     false,
		wantOddY:  "",
		wantEvenY: "",
	}}

	for _, test := range tests {
		// Decompress the test odd y coordinate for the given test x coordinate
		// and ensure the returned validity flag matches the expected result.
		var oddY fieldVal
		fx := new(fieldVal).SetHex(test.x)
		valid := decompressY(fx, true, &oddY)
		if valid != test.valid {
			t.Errorf("%s: unexpected valid flag -- got: %v, want: %v",
				test.name, valid, test.valid)
			continue
		}

		// Decompress the test even y coordinate for the given test x coordinate
		// and ensure the returned validity flag matches the expected result.
		var evenY fieldVal
		valid = decompressY(fx, false, &evenY)
		if valid != test.valid {
			t.Errorf("%s: unexpected valid flag -- got: %v, want: %v",
				test.name, valid, test.valid)
			continue
		}

		// Skip checks related to the y coordinate when there isn't one.
		if !valid {
			continue
		}

		// Ensure the decompressed odd Y coordinate is the expected value.
		oddY.Normalize()
		wantOddY := new(fieldVal).SetHex(test.wantOddY)
		if !wantOddY.Equals(&oddY) {
			t.Errorf("%s: mismatched odd y\ngot: %v, want: %v", test.name,
				oddY, wantOddY)
			continue
		}

		// Ensure the decompressed even Y coordinate is the expected value.
		evenY.Normalize()
		wantEvenY := new(fieldVal).SetHex(test.wantEvenY)
		if !wantEvenY.Equals(&evenY) {
			t.Errorf("%s: mismatched even y\ngot: %v, want: %v", test.name,
				evenY, wantEvenY)
			continue
		}

		// Ensure the decompressed odd y coordinate is actually odd.
		if !oddY.IsOdd() {
			t.Errorf("%s: odd y coordinate is even", test.name)
			continue
		}

		// Ensure the decompressed even y coordinate is actually even.
		if evenY.IsOdd() {
			t.Errorf("%s: even y coordinate is odd", test.name)
			continue
		}
	}
}

// TestDecompressYRandom ensures that decompressY works as expected with
// randomly-generated x coordinates.
func TestDecompressYRandom(t *testing.T) {
	// Use a unique random seed each test instance and log it if the tests fail.
	seed := time.Now().Unix()
	rng := rand.New(rand.NewSource(seed))
	defer func(t *testing.T, seed int64) {
		if t.Failed() {
			t.Logf("random seed: %d", seed)
		}
	}(t, seed)

	for i := 0; i < 100; i++ {
		origX := randFieldVal(t, rng)

		// Calculate both corresponding y coordinates for the random x when it
		// is a valid coordinate.
		var oddY, evenY fieldVal
		x := new(fieldVal).Set(origX)
		oddSuccess := decompressY(x, true, &oddY)
		evenSuccess := decompressY(x, false, &evenY)

		// Ensure that the decompression success matches for both the even and
		// odd cases depending on whether or not x is a valid coordinate.
		if oddSuccess != evenSuccess {
			t.Fatalf("mismatched decompress success for x = %v -- odd: %v, "+
				"even: %v", x, oddSuccess, evenSuccess)
		}
		if !oddSuccess {
			continue
		}

		// Ensure the x coordinate was not changed.
		if !x.Equals(origX) {
			t.Fatalf("x coordinate changed -- orig: %v, changed: %v", origX, x)
		}

		// Ensure that the resulting y coordinates match their respective
		// expected oddness.
		oddY.Normalize()
		evenY.Normalize()
		if !oddY.IsOdd() {
			t.Fatalf("requested odd y is even for x = %v", x)
		}
		if evenY.IsOdd() {
			t.Fatalf("requested even y is odd for x = %v", x)
		}

		// Ensure that the resulting x and y coordinates are actually on the
		// curve for both cases.
		if !isOnCurve(x, &oddY) {
			t.Fatalf("(%v, %v) is not a valid point", x, oddY)
		}
		if !isOnCurve(x, &evenY) {
			t.Fatalf("(%v, %v) is not a valid point", x, evenY)
		}
	}
}
