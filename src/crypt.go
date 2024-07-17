// Copyright (c) 2024 Nic Neate
//
// Dummy cryptography functions for the MPC protocol, which are used to encrypt and decrypt
// the outputs of the circuit gates.
//
// Rather than using actual encryption, these functions are design to provide readable output
// so it is easy to understand a printed version of the garbled gate.
package mpc

import (
	"fmt"
	"log"
	"strconv"
)

type Key struct {
	secret  string
	pointer int
}

func generate_key(input_idx int, input_val bool, pointer int) Key {
	secret := fmt.Sprintf("key%d%v", input_idx, input_val)
	return Key{secret: secret, pointer: pointer}
}

func encrypt(value bool, key_a, key_b Key) string {
	return fmt.Sprintf("%v-encrypted-by-%s-%s", value, key_a.secret, key_b.secret)
}

func decrypt(value string, _, _ Key) bool {
	bool_val, err := strconv.ParseBool(value[:1])
	if err != nil {
		log.Fatal(err)
	}
	return bool_val
}

func key_to_string(key Key) string {
	return fmt.Sprintf("%s,%d", key.secret, key.pointer)
}

func key_from_string(key_str string) Key {
	pointer, err := strconv.Atoi(key_str[len(key_str)-1:])
	if err != nil {
		log.Fatal(err)
	}
	return Key{secret: key_str[:len(key_str)-2], pointer: pointer}
}
