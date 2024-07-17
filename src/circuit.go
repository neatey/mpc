// Copyright (c) 2024 Nic Neate
//
// Types representing boolean logic circuit as used in Yao's garbled circuit
// protocol
package mpc

// A gate is a single boolean logic gate (AND, OR, XOR, NOT etc.) that takes one
// or two boolean inputs and produces a single boolean output.
//
// For now, we only implement a single type: the AND gate.
// EXTENSION: Implement all of the fundamental boolean logic gates.
type Gate struct {
	gate_type string // Always "AND"
}

func (g *Gate) Output(input_a bool, input_b bool) bool {
	// For now, we only implement a single AND gate
	return input_a && input_b
}

// Circuit is a representation of a boolean logic circuit that can be garbled
// and evaluated using Yao's protocol. It is a collection of one or more gates
// and wires that connect them.
//   - Input wires provide the input to the circuit and each is connected to
//     one or more gates.
//   - The output wire of each gate is connected to zero or more other gates. If
//     an output wire is not connect to another gate, it is an output of the
//     circuit.
//
// For now, we only implement a single AND gate.
// EXTENSION: Support arbitrary circuits with multiple gates and wires.
type Circuit struct {
	gate Gate
}

// A garbled gate is a matrix of encrypted outputs for each possible
// combination of input values.
type GarbledGate [2][2]string

// A garbled circuit is a collection of garbled gates. For now, we only support
// a single gate.
type GarbledCircuit struct {
	garbled_gate GarbledGate
}
