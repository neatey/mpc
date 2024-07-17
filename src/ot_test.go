// Copyright (c) 2024 Nic Neate
package mpc

import (
	"fmt"
	"testing"
)

func TestObliviousTransfer(t *testing.T) {
	fmt.Println("")
	fmt.Println("Test case: Perform Oblivous Transfer of one of two messages")

	ot_sender := OtSender{message0: "message0", message1: "message1"}
	ot_receiver := OtReceiver{choice: false}
	fmt.Printf("Message 0: \"%s\" %v\n", ot_sender.message0, []byte(ot_sender.message0))
	fmt.Printf("Message 1: \"%s\" %v\n", ot_sender.message1, []byte(ot_sender.message1))

	message := PerformObliviousTransfer(ot_sender, ot_receiver)
	fmt.Printf("Transferred message: \"%s\" %v\n", message, []byte(message))

	if message != "message0" {
		t.Fatalf("Bad message received: %s %v", message, []byte(message))
	}
}
