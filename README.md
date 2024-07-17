# mpc

Toy project to illustrate secure multi-party computation (MPC) using Go.

This project provides a very simple example implementation in Go of two-party computation (2PC) using Yao's garbled circuit protocol with oblivious transfer.

## Getting started

Simply run `go test` in the [src](./src) directory. [2pc_test.go](./src/2pc_test.go) provides a walked example of a trivial 2PC to demonstrate the protocol - take a look and read the comments to understand how the protocol works.

## Background on secure multi-party computation (MPC)

Secure multi-party computation allows multiple parties to collaborate and compute something without sharing their inputs with each other, and without the need for a trusted third party.

**Example 1:** Two millionaires want to know who is richer, but don't want to tell each other their net worth.

**Example 2:** To provide a trusted signature on an important image, a private key is divided between two developers and both need to provide their parts to the signing function but without sharing their private key part with anyone.

### Why aren't we all using this already?

The performance is atrocious. Any interesting computation requires 1000s of gates, and the garbler has to compute and encrypt the result of every possible input on every gate. In total, the protocol requires four polynomial-time algorithms. Only recently have optimizations made it potentially viable for practical use.

Take example 2 above, image signing.
-  An optimised GC circuit for SHA256 hashing contains around 120k gates.
-  ed25519 2-part key generation takes ~1.5mins to run.
-  ed25519 signing takes a further ~1.5mins to run.
-  That's much too long to use in e.g. TLS. It's a viable solution for having multiple developers sign an image. But it's not any more secure than just attaching multiple signatures to the image, with each developer signing it in the traditional way.

### Yao's Garbled Circuit protocol

Yao’s Garbled Circuit protocol is an MPC protocol for two parties (i.e. 2PC). Here's a quick explanation of how it works:

1. **Boolean Circuit Representation:** The function to be computed is represented as a Boolean circuit, consisting of basic gates like AND, OR, and XOR.
2. **Garbled Circuit Creation:** One party, called the “garbler,” creates a “garbled” version of the circuit. This involves encrypting the inputs and outputs of each gate in such a way that the evaluator can compute the function without learning the intermediate values. Cruicially, that means using a unique encryption key for each possible value of each input. 
3. **Input Encryption:** The garbler encrypts their own input(s) and sends the garbled circuit along with these encrypted input(s) to the other party, called the “evaluator.”
4. **Oblivious Transfer:** The evaluator needs to encrypt their inputs as well. This is tricky because the only the garbler knows the right encryption key(s) for the value(s) of evaluator's input(s), but they don't know those value(s). The solution is Oblivous Transfer, a technique which allows the evaluator to request the right key(s) without the garbler knowing which one(s) it sent.
5. **Circuit Evaluation:** The evaluator uses the garbled circuit and the encryption keys to compute the function. They can decrypt the final output but cannot learn anything about the intermediate values or the garbler’s inputs.
6. **Output:** Finally, the evaluator can reveal the final output to the garbler.

This protocol ensures that neither party learns anything about the other’s inputs beyond what can be inferred from the output.

Didn't make sense? Need more detail? Try these references:

-  A relatively brief, plain english explanation with some nice diagrams: [Explaining Yao's Garbled Circuits | cronokirby.com](https://cronokirby.com/posts/2022/05/explaining-yaos-garbled-circuits/)
-  A couple of more academic explanations that go into further details, extensions and optimizations: [A Gentle Introduction to Yao’s Garbled Circuits | Yakoubov | Boston University](https://web.mit.edu/sonka89/www/papers/2017ygc.pdf) and [Yao’s Garbled Circuits: Recent Directions and Implementations | Snyder | University of Illinois at Chicago](https://www.peteresnyder.com/static/papers/Peter_Snyder_-_Garbled_Circuits_WCP_2_column.pdf). 

### 2-1 Oblivious Transfer

There's a good explanation of the Oblibvious Transfer process here: [Oblivious Transfer | wikiwand](https://www.wikiwand.com/en/Oblivious_transfer) - but note that there are typos in some of the RSA functions.

## Overview of this repo

-  [2pc_test.go](./src/2pc_test.go) is a worked example of 2PC using Garbled Circuits and Oblivious Transfer. Run the test using `go test`. Follow what it's doing by reading the comments in the code, and stepping into the (also commented) functions it calls for the different steps in the protocol.
-  [circuit.go](./src/circuit.go) provides the data structures describing the circuit itself. Currently just a trivial circuit containing a single AND gate is supported.
-  [yao.go](./src/yao.go) implements the Garbled Circuit protocol itself (in a simplistic and unoptimized way, to make it easy to read and understand the protocol).
-  [ot.go](./src/ot.go) implements 2-1 Oblivious Transfer using RSA encryption between the sender and receiver. There's also a test in [ot_test.go](./src/ot_test.go).
-  [crypt.go](./src/crypt.go) is a utililty module providing some dummy cryptographic functions to use in the protocol.

## Making this useful

There are plenty of good open source implementations out there already which are good enough quality to be used in real life applications. This one is a good example: [mpc | Markku Rossi](https://github.com/markkurossi/mpc/blob/master/README.md). The README include instructions to run the ed25519 signing example given above.

To turn my (this) repo into something actually useful, we'd need to make it look more like Markku's one:

-  Implement an extenion to be secure against one party misusing the protocol.
-  Add support for arbitrary circuits with boolean logic gates of any type.
-  Run the garbler and evaluator as separate binaries communicating over a network and not sharing any other data.
-  Use real encryption.
-  Write a complier to convert ~normal code into a boolean circuit representation.
-  Optimise everything.