// Implementation of the Yao's Garbled Circuit Protocol. We define two types, Garbler and Evaluator,
// which proform the roles of the two parties in the protocol.
//
// Refer to 2pc_test.go to see how these types are used to run the protocol.
package mpc

type Garbler struct {
	input bool
	keys  [2]map[bool]Key
}

func (g *Garbler) Garble(circuit Circuit) GarbledCircuit {
	// Generate an encryption key for each possible value of each input
	for input_idx := 0; input_idx < 2; input_idx++ {
		g.keys[input_idx] = make(map[bool]Key)
		for _, input_val := range []bool{true, false} {
			g.keys[input_idx][input_val] = generate_key(input_idx, input_val)
		}
	}

	// Encrypt each entry in the output matrix
	var encrypted_output GarbledCircuit
	for a, input_a := range []bool{true, false} {
		for b, input_b := range []bool{true, false} {
			output := circuit.Output(input_a, input_b)
			encrypted_output[a][b] = encrypt(output, g.keys[0][input_a], g.keys[1][input_b])
		}
	}

	// Randomize the 'a' index in the encrypted output matrix
	// TODO: Implement this

	// Randomize the 'b' index in the encrypted output matrix
	// TODO: Implement this

	return encrypted_output
}

func (g *Garbler) TransferKeys(input_b bool) (Key, Key) {
	// Return the two keys that the Evaluator must use to evaluate the circuit.
	// TODO: Implement Oblivious Transfer so that the Garbler does not learn the
	// value of input b.
	key_a := g.keys[0][g.input]
	key_b := g.keys[1][input_b]
	return key_a, key_b
}

type Evaluator struct {
	input bool
}

func (e *Evaluator) Evaluate(garbled_circuit GarbledCircuit, key_a, key_b Key) bool {
	// To evaluate, look-up the encrypted output in the garbled circuit matrix that corresponds
	// to the row/column of each key, then decrypt using the two keys
	encrypted_output := garbled_circuit[key_a.pointer][key_b.pointer]
	return decrypt(encrypted_output, key_a, key_b)
}
