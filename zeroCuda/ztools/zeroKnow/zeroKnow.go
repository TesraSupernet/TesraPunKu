package zeroKnow

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/arnaucube/go-snark"
	"github.com/arnaucube/go-snark/circuitcompiler"
	"github.com/arnaucube/go-snark/groth16"
	"github.com/arnaucube/go-snark/r1csqap"
	"math/big"
)

type ZkData struct {
	CompileJson   []byte
	CircuitFun    []byte
	PublicInputs  []byte
	PrivateInputs []byte
	TrustedSetup  []byte
	Proofs        []byte
}

func (z *ZkData) CreatZeroProofs() error {

	compileJsonDate, err := CompileCircuit(z.CircuitFun, z.PublicInputs, z.PrivateInputs)
	if err != nil {
		//fmt.Println("err=", err)
		return err
	}
	z.CompileJson = compileJsonDate
	//fmt.Println("compileJsonDate=", string(compileJsonDate))

	z.TrustedSetup, err = Groth16TrustedSetup(compileJsonDate, z.PublicInputs, z.PrivateInputs)
	if err != nil {
		//fmt.Println("err=", err)
		return err
	}
	//fmt.Println("trusted=", string(trusted))

	z.Proofs, err = Groth16GenerateProofs(compileJsonDate, z.TrustedSetup, z.PublicInputs, z.PrivateInputs)
	if err != nil {
		//fmt.Println("err=", err)
		return err
	}
	return nil
}

func ZeroVerify(proofs []byte, trustedSetup []byte, publicInputs []byte) (bool, error) {
	ok, err := Groth16VerifyProofs(proofs, trustedSetup, publicInputs)
	if err != nil {
		//fmt.Println("err=", err)
		return false, err
	}
	if ok == false {
		//fmt.Println("verify is error")
		return false, nil
	}
	return true, nil
}

func Groth16GenerateProofs(compileCircuitBytes []byte, trustedSetupBytes []byte, publicInput []byte, privateInput []byte) ([]byte, error) {
	var circuit circuitcompiler.Circuit
	err := json.Unmarshal(compileCircuitBytes, &circuit)
	if err != nil {
		return nil, err
	}
	// open trustedsetup.json
	var trustedsetup groth16.Setup
	err = json.Unmarshal(trustedSetupBytes, &trustedsetup)
	if err != nil {
		return nil, err
	}
	// parse inputs from inputsFile
	var inputs circuitcompiler.Inputs
	err = json.Unmarshal(privateInput, &inputs.Private)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(publicInput, &inputs.Public)
	if err != nil {
		return nil, err
	}
	// calculate wittness
	w, err := circuit.CalculateWitness(inputs.Private, inputs.Public)
	if err != nil {
		return nil, err
	}
	//fmt.Println("witness", w)

	// flat code to R1CS
	a := circuit.R1CS.A
	b := circuit.R1CS.B
	c := circuit.R1CS.C
	// R1CS to QAP
	alphas, betas, gammas, _ := groth16.Utils.PF.R1CSToQAP(a, b, c)
	_, _, _, px := groth16.Utils.PF.CombinePolynomials(w, alphas, betas, gammas)
	groth16.Utils.PF.DivisorPolynomial(px, trustedsetup.Pk.Z)

	proof, err := groth16.GenerateProofs(circuit, trustedsetup.Pk, w, px)
	if err != nil {
		return nil, err
	}
	// store Proofs to json
	jsonData, err := json.Marshal(proof)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func Groth16VerifyProofs(proofBytes []byte, trustedSetupBytes []byte, publicInputs []byte) (bool, error) {
	// open Proofs.json
	var proof groth16.Proof
	err := json.Unmarshal(proofBytes, &proof)
	if err != nil {
		return false, err
	}
	var trustedsetup groth16.Setup
	err = json.Unmarshal(trustedSetupBytes, &trustedsetup)
	if err != nil {
		return false, err
	}
	// read publicInputs file

	var publicSignals []*big.Int
	err = json.Unmarshal(publicInputs, &publicSignals)
	if err != nil {
		return false, err
	}
	verified := groth16.VerifyProof(trustedsetup.Vk, proof, publicSignals, true)
	if !verified {
		return false, nil
	}
	return true, nil
}

func Groth16TrustedSetup(compileCircuitBytes []byte, publicInput []byte, privateInput []byte) ([]byte, error) {
	var circuit circuitcompiler.Circuit
	err := json.Unmarshal(compileCircuitBytes, &circuit)
	if err != nil {
		return nil, err
	}

	// parse inputs from inputsFile
	var inputs circuitcompiler.Inputs
	err = json.Unmarshal(privateInput, &inputs.Private)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(publicInput, &inputs.Public)
	if err != nil {
		return nil, err
	}

	// calculate wittness
	w, err := circuit.CalculateWitness(inputs.Private, inputs.Public)
	if err != nil {
		return nil, err
	}

	// R1CS to QAP
	alphas, betas, gammas, _ := snark.Utils.PF.R1CSToQAP(circuit.R1CS.A, circuit.R1CS.B, circuit.R1CS.C)

	// calculate trusted setup
	setup, err := groth16.GenerateTrustedSetup(len(w), circuit, alphas, betas, gammas)
	if err != nil {
		return nil, err
	}
	//fmt.Println("\nt:", setup.Toxic.T)

	// remove setup.Toxic
	var tsetup groth16.Setup
	tsetup.Pk = setup.Pk
	tsetup.Vk = setup.Vk

	// store setup to json
	jsonData, err := json.Marshal(tsetup)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func CompileCircuit(circuitBytes, publicInputs, privateInputs []byte) ([]byte, error) {
	//var (
	//	compileJsonDate []byte
	//	pxJsonDate      []byte
	//)

	// read circuit file

	// parse circuit code
	parser := circuitcompiler.NewParser(bytes.NewReader(circuitBytes))
	circuit, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	//fmt.Println("\ncircuit data:", circuit)

	// parse inputs from inputsFile
	var inputs circuitcompiler.Inputs
	err = json.Unmarshal(privateInputs, &inputs.Private)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(publicInputs, &inputs.Public)
	if err != nil {
		return nil, err
	}

	// calculate wittness
	w, err := circuit.CalculateWitness(inputs.Private, inputs.Public)
	if err != nil {
		return nil, err
	}
	//fmt.Println("\nwitness", w)

	// flat code to R1CS
	//fmt.Println("\ngenerating R1CS from flat code")
	a, b, c := circuit.GenerateR1CS()
	// R1CS to QAP
	alphas, betas, gammas, zx := snark.Utils.PF.R1CSToQAP(a, b, c)

	ax, bx, cx, px := snark.Utils.PF.CombinePolynomials(w, alphas, betas, gammas)

	hx := snark.Utils.PF.DivisorPolynomial(px, zx)

	// hx==px/zx so px==hx*zx
	// assert.Equal(t, px, snark.Utils.PF.Mul(hx, zx))
	if !r1csqap.BigArraysEqual(px, snark.Utils.PF.Mul(hx, zx)) {
		return nil, errors.New("px != hx*zx")
	}

	// p(x) = a(x) * b(x) - c(x) == h(x) * z(x)
	abc := snark.Utils.PF.Sub(snark.Utils.PF.Mul(ax, bx), cx)
	// assert.Equal(t, abc, px)
	if !r1csqap.BigArraysEqual(abc, px) {
		return nil, errors.New("abc != px")
	}
	hz := snark.Utils.PF.Mul(hx, zx)
	if !r1csqap.BigArraysEqual(abc, hz) {
		return nil, errors.New("abc != hz")
	}
	// assert.Equal(t, abc, hz)

	div, rem := snark.Utils.PF.Div(px, zx)
	if !r1csqap.BigArraysEqual(hx, div) {
		return nil, errors.New("hx != div")
	}
	// assert.Equal(t, hx, div)
	// assert.Equal(t, rem, r1csqap.ArrayOfBigZeros(4))
	for _, r := range rem {
		if !bytes.Equal(r.Bytes(), big.NewInt(int64(0)).Bytes()) {
			panic(errors.New("error:error:  px/zx rem not equal to zeros"))
		}
	}

	// store circuit to json
	jsonDate, err := json.Marshal(circuit)
	if err != nil {
		return nil, err
	}

	return jsonDate, nil
}
