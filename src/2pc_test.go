package mpc

import (
	"testing"
)

func TestDemonstrate2pc(t *testing.T) {
	// Define the circuit for the muilt-party computation we are going to run.
	// For now, we only support a single AND gate.
	// EXTENSION: Define a JSON form for multi-gate circuits and read it in from
	// an input file.
	circuit := Circuit{}
	input_a := true
	input_b := false

	// Create the two parties involved in this computation: the Garbler and the
	// Evaluator. Give each of them their input value for the computation. They
	// do not share these values with each other.
	garbler := Garbler{input: input_a}
	evaluator := Evaluator{input: input_b}

	// Step 1: Garbler creates the garbled circuit.
	garbled_circuit := garbler.Garble(circuit)

	// Step 2: Garbler transfers the encryption keys for each of the two inputs
	// to the Evaluator.
	// EXTENSION: Implement Oblivious Transfer so that the Garbler does not learn
	// the value of the Evaluator's input.
	key_a, key_b := garbler.TransferKeys(evaluator.input)

	// Step 3: Evaluator computes the result, which can then be returned to the
	// Garbler as well.
	output := evaluator.Evaluate(garbled_circuit, key_a, key_b)

	if output != circuit.Output(input_a, input_b) {
		t.Fatalf("Circuit is expected to evalulate to %v", !output)
	}
}
