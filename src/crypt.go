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

func generate_key(input_idx int, input_val bool) Key {
	secret := fmt.Sprintf("key-%d-%v", input_idx, input_val)
	return Key{secret: secret, pointer: input_idx}
}

func encrypt(value bool, key_a, key_b Key) string {
	return fmt.Sprintf("%v-encrypted-%s-%s", value, key_a.secret, key_b.secret)
}

func decrypt(value string, _, _ Key) bool {
	bool_val, err := strconv.ParseBool(value[:1])
	if err != nil {
		log.Fatal(err)
	}
	return bool_val
}
