// Copyright (c) 2024 Nic Neate
package mpc

import (
	"fmt"
	"testing"
)

func TestDemonstrate2pc(t *testing.T) {
	// Define the circuit for the muilt-party computation we are going to run.
	// For now, we only support a single AND gate.
	// EXTENSION: Define a JSON form for multi-gate circuits and read it in from
	// an input file.
	circuit := Circuit{Gate{"AND"}}
	input_a := false
	input_b := true

	fmt.Println("Evaluating circuit:", circuit)
	fmt.Println("Input A:", input_a)
	fmt.Println("Input B:", input_b)

	// Create the two parties involved in this computation: the Garbler and the
	// Evaluator. Give each of them their input value for the computation. They
	// do not share these values with each other.
	garbler := Garbler{input: input_a}
	evaluator := Evaluator{input: input_b}

	// Step 1: Garbler creates the garbled circuit: a matrix of the outputs of
	// each gate, each encrypted with a unique pair of keys corresponding to the
	// inputs.
	garbled_circuit := garbler.GarbleCircuit(circuit)

	fmt.Println("Garbled circuit:", garbled_circuit)

	// Step 2: Garbler transfers the encryption keys for each of the two inputs
	// to the Evaluator.
	//
	// The Garbler knows what input_a is, and can simply provide the encryption key.
	// However, the Garbler cannot know input_b, so we use 1-2 Oblivious Transfer to
	// allow the Evaluator to select the encryption key corresponding to input_b.
	key_a := garbler.GetKeyA()
	ot_sender := garbler.GetKeyBOtSender()
	ot_receiver := evaluator.GetKeyBOtReceiver()
	key_b := PerformObliviousTransfer(ot_sender, ot_receiver)

	fmt.Println("Transferred key A:", key_a)
	fmt.Println("Transferred key B:", key_b)

	// Step 3: Evaluator uses the keys to computes the result of the circuit, which
	// can then be returned to the Garbler as well.
	output := evaluator.Evaluate(garbled_circuit, key_a, key_b)

	fmt.Println("Output:", output)

	if output != circuit.gate.Output(input_a, input_b) {
		t.Fatalf("Circuit is expected to evalulate to %v", !output)
	}
}
