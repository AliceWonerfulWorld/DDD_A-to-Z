package brainfuck

import (
	"errors"
	"fmt"
)

const (
	tapeSize       = 30000
	maxInstruction = 1000000
)

var (
	ErrBracketMismatch = errors.New("brainfuck brackets are not balanced")
	ErrPointerBounds   = errors.New("brainfuck pointer moved out of bounds")
	ErrStepLimit       = errors.New("brainfuck execution step limit exceeded")
)

func ExecuteBF(code string, input string) (string, error) {
	jumps, err := buildJumpTable(code)
	if err != nil {
		return "", err
	}

	tape := make([]byte, tapeSize)
	inputBytes := []byte(input)
	output := make([]byte, 0, len(inputBytes))
	var pc, ptr, inputPos, steps int

	for pc < len(code) {
		steps++
		if steps > maxInstruction {
			return "", ErrStepLimit
		}

		switch code[pc] {
		case '+':
			tape[ptr]++
		case '-':
			tape[ptr]--
		case '>':
			ptr++
			if ptr >= tapeSize {
				return "", ErrPointerBounds
			}
		case '<':
			ptr--
			if ptr < 0 {
				return "", ErrPointerBounds
			}
		case '.':
			output = append(output, tape[ptr])
		case ',':
			if inputPos < len(inputBytes) {
				tape[ptr] = inputBytes[inputPos]
				inputPos++
			} else {
				tape[ptr] = 0
			}
		case '[':
			if tape[ptr] == 0 {
				pc = jumps[pc]
			}
		case ']':
			if tape[ptr] != 0 {
				pc = jumps[pc]
			}
		}

		pc++
	}

	return string(output), nil
}

func buildJumpTable(code string) (map[int]int, error) {
	jumps := map[int]int{}
	stack := make([]int, 0)

	for pc := range code {
		switch code[pc] {
		case '[':
			stack = append(stack, pc)
		case ']':
			if len(stack) == 0 {
				return nil, fmt.Errorf("%w: unexpected ] at %d", ErrBracketMismatch, pc)
			}
			open := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			jumps[open] = pc
			jumps[pc] = open
		}
	}
	if len(stack) > 0 {
		return nil, fmt.Errorf("%w: missing ] for [ at %d", ErrBracketMismatch, stack[len(stack)-1])
	}

	return jumps, nil
}
