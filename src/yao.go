// Copyright (c) 2024 Nic Neate
//
// Implementation of the Yao's Garbled Circuit Protocol. We define two types, Garbler and Evaluator,
// which proform the roles of the two parties in the protocol.
//
// Refer to 2pc_test.go to see how these types are used to run the protocol.
package mpc

import "math/rand"

type Garbler struct {
	input bool
	keys  [2]map[bool]Key
}

func (g *Garbler) GarbleCircuit(circuit Circuit) GarbledCircuit {
	// For now we only support a single gate, so we just have to garble that.
	garbled_gate := g.GarbleGate(circuit.gate)
	return GarbledCircuit{garbled_gate}
}

func (g *Garbler) GarbleGate(gate Gate) GarbledGate {
	// Generate an encryption key for each possible value of each input
	for input_idx := 0; input_idx < 2; input_idx++ {
		g.keys[input_idx] = make(map[bool]Key)
		for input_val_idx, input_val := range []bool{true, false} {
			g.keys[input_idx][input_val] = generate_key(input_idx, input_val, input_val_idx)
		}
	}

	// Calculate the output of the gate for each combination of input values, and generate an
	// encrypted matrix of those values: the garbed gate.
	var garbled_gate GarbledGate
	for a, input_a := range []bool{true, false} {
		for b, input_b := range []bool{true, false} {
			output := gate.Output(input_a, input_b)
			garbled_gate[a][b] = encrypt(output, g.keys[0][input_a], g.keys[1][input_b])
		}
	}

	// Randomize the 'a' index in the encrypted output matrix, and update the key pointers
	// accordingly
	if rand.Intn(2) == 1 {
		garbled_gate[0], garbled_gate[1] = garbled_gate[1], garbled_gate[0]
		key_true, key_false := g.keys[0][true], g.keys[0][false]
		key_true.pointer, key_false.pointer = key_false.pointer, key_true.pointer
		g.keys[0][true], g.keys[0][false] = key_true, key_false
	}

	// Randomize the 'b' index in the encrypted output matrix, and update the key pointers
	// accordingly
	if rand.Intn(2) == 1 {
		garbled_gate[0][0], garbled_gate[0][1] = garbled_gate[0][1], garbled_gate[0][0]
		garbled_gate[1][0], garbled_gate[1][1] = garbled_gate[1][1], garbled_gate[1][0]
		key_true, key_false := g.keys[1][true], g.keys[1][false]
		key_true.pointer, key_false.pointer = key_false.pointer, key_true.pointer
		g.keys[1][true], g.keys[1][false] = key_true, key_false
	}

	return garbled_gate
}

func (g *Garbler) GetKeyA() string {
	// Return the key for input a that the Evaluator must use to evaluate the circuit.
	key := g.keys[0][g.input]
	return key_to_string(key)
}

func (g *Garbler) GetKeyBOtSender() OtSender {
	// Create an Oblivious Transfer sender containing two messages: the string from of
	// each of the keys for input b
	key_true := g.keys[1][true]
	key_false := g.keys[1][false]
	return OtSender{message0: key_to_string(key_false), message1: key_to_string(key_true)}
}

type Evaluator struct {
	input bool
}

func (e *Evaluator) GetKeyBOtReceiver() OtReceiver {
	// Create an Oblivious Transfer receiver that will select the appropriate key for
	// input b from the Garbler.
	return OtReceiver{choice: e.input}
}

func (e *Evaluator) Evaluate(garbled_circuit GarbledCircuit, key_a_str, key_b_str string) bool {
	// To evaluate, look-up the encrypted output in the garbled circuit matrix that corresponds
	// to the row/column of each key, then decrypt using the two keys
	key_a := key_from_string(key_a_str)
	key_b := key_from_string(key_b_str)
	encrypted_output := garbled_circuit.garbled_gate[key_a.pointer][key_b.pointer]
	return decrypt(encrypted_output, key_a, key_b)
}
