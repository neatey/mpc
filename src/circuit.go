package mpc

// Circuit is a representation of a boolean logic circuit that can be garbled
// and evaluated using Yao's protocol.
type Circuit struct {
}

func (c *Circuit) Output(input_a bool, input_b bool) bool {
	// For now, we only implement a single AND gate
	return input_a && input_b
}

// The garbled circuit is a matrix of encrypted outputs for each possible
// combination of input values.
type GarbledCircuit [2][2]string
