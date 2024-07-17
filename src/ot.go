// Copyright (c) 2024 Nic Neate
//
// A simple implementatino of 1-2 Oblivious Transfer (OT) for use with the Garbled
// Circuit protocol.
package mpc

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

// Perform a 1-2 Oblivious Transfer between a sender and a receiver. This is an illustrative
// implementation of the protocol which doesn't bother with actual encryption or sending the
// messages over a network. It demonstrates the protocol as follows:
//   - The sender passed in has two messages and needs to send exactly one of them to the receiver.
//   - They receiver has a bit that determines which message they want to receive. They need to
//     keep that choice secret from the receiver.
//   - This function keeps all secret information within the sender and receiver types. The
//     local variables exposed in this function are all public information.
func PerformObliviousTransfer(sender OtSender, receiver OtReceiver) string {
	// The sender generates an RSA key pair and shares the public key with the receiver
	public_key := sender.GetPublicKey()

	// The sender generates two random values x0 and x1 that the receiver will use to
	// choose between the messages.
	x0, x1 := sender.GetRandomValuesX0X1()

	// The receiver selects xb - the random value corresponding to the message they went.
	// They then generate a their own random value k and use it to blind xb, so that the
	// sender does not know which they have chosen.
	v := receiver.BlindXb(x0, x1, public_key)

	// The sender uses the value v to encrypt both messages, and sends them to the receiver.
	// The sender never knows which one the receiver has chosen. The receiver can only decrypt
	// the one they have chosen.
	encrypted_message0, encrypted_message1 := sender.GetEncryptedMessages(v)
	fmt.Println("-- OT Encypted message 0:", encrypted_message0)
	fmt.Println("-- OT Encypted message 1:", encrypted_message1)

	// Finally, the receiver decrypts their selected message.
	return receiver.DecryptSelectedMessage(encrypted_message0, encrypted_message1, public_key.modulus_n)
}

// Structs to define a simple (just using ints) RSA key pair
type RsaKeyPair struct {
	public_key  RsaPublicKey
	private_key RsaPrivateKey
}

type RsaPublicKey struct {
	modulus_n         int64
	public_exponent_e int64
}

type RsaPrivateKey struct {
	private_exponent_d int64
}

// Generate a new RSA key pair. Because this is a toy implementation designed
// to illustrate the protocol, we use a trivial fixed key pair that allows us
// to work with int types:
// p = 61, q = 53
// N = p * q = 3233
// e = 17
// e * d = 1 mod (p-1)(q-1) = 1 mod 3120
// d = 2753
func generate_rsa_key_pair() RsaKeyPair {
	return RsaKeyPair{
		public_key: RsaPublicKey{
			modulus_n:         3233,
			public_exponent_e: 17,
		},
		private_key: RsaPrivateKey{
			private_exponent_d: 2753,
		},
	}
}

func pow_mod(a, b, m int64) int64 {
	out := big.NewInt(0)
	return out.Exp(big.NewInt(int64(a)), big.NewInt(int64(b)), big.NewInt(int64(m))).Int64()
}

// Compute a mod b, but always return a positive number
func mod(a float64, b int64) int64 {
	m := int64(math.Mod(a, float64(b)))
	if m < 0 {
		m += b
	}
	return m
}

// Encrypt a string byte by byte
func encrypt_message(message string, k int64, modulus_n int64) []int64 {
	var encrypted_message []int64
	for _, c := range message {
		encrypted_message = append(encrypted_message, mod(float64(c)+float64(k), modulus_n))
	}
	return encrypted_message
}

// Decrypt a string by reversing the encrypt() function above
func decrypt_message(encrypted_message []int64, k int64, modulus_n int64) string {
	var decrypted_runes []rune
	for _, ec := range encrypted_message {
		decrypted_runes = append(decrypted_runes, rune(mod(float64(int64(ec)-k), modulus_n)))
	}
	return string(decrypted_runes)
}

type OtSender struct {
	message0     string
	message1     string
	rsa_key_pair RsaKeyPair
	x0           int64
	x1           int64
}

func (s *OtSender) GetPublicKey() RsaPublicKey {
	s.rsa_key_pair = generate_rsa_key_pair()
	return s.rsa_key_pair.public_key
}

// Generate and return two random values.
func (s *OtSender) GetRandomValuesX0X1() (int64, int64) {
	s.x0 = rand.Int63n(100)
	s.x1 = rand.Int63n(100)
	return s.x0, s.x1
}

func (s *OtSender) compute_k0_k1(v int64) (int64, int64) {
	// Compute k0 and k1 using the RSA private key
	// xb + k = decrypted(v) = v^d mod N
	// kb = (v^d - xb) mod N
	xb_plus_kb := pow_mod(v, s.rsa_key_pair.private_key.private_exponent_d, s.rsa_key_pair.public_key.modulus_n)

	k0 := mod(float64(xb_plus_kb-s.x0), s.rsa_key_pair.public_key.modulus_n)
	k1 := mod(float64(xb_plus_kb-s.x1), s.rsa_key_pair.public_key.modulus_n)
	return int64(k0), int64(k1)
}

func (s *OtSender) GetEncryptedMessages(v int64) ([]int64, []int64) {
	k0, k1 := s.compute_k0_k1(v)

	encrypted_message0 := encrypt_message(s.message0, k0, s.rsa_key_pair.public_key.modulus_n)
	encrypted_message1 := encrypt_message(s.message1, k1, s.rsa_key_pair.public_key.modulus_n)

	return encrypted_message0, encrypted_message1
}

type OtReceiver struct {
	choice bool
	k      int64
}

func (r *OtReceiver) BlindXb(x0, x1 int64, public_key RsaPublicKey) int64 {
	// Choose xb and a random k
	xb := x0
	if r.choice {
		xb = x1
	}
	r.k = rand.Int63n(100)

	// RSA encrypt xb + k using (xb + k) ^ e mod N
	blind_xb := pow_mod(xb+r.k, public_key.public_exponent_e, public_key.modulus_n)

	return blind_xb
}

func (r *OtReceiver) DecryptSelectedMessage(encrypted_message0, encrypted_message1 []int64, modulus_n int64) string {
	encrypted_message := encrypted_message0
	if r.choice {
		encrypted_message = encrypted_message1
	}
	return decrypt_message(encrypted_message, r.k, modulus_n)
}
