#!/bin/sh

echo 'draw a dog' |
./cmd/txt2llm2png/txt2llm2png \
	--model x/flux2-klein:4b-fp4 \
	--width 128 \
	--height 128 \
	--steps 4 \
	--seed 101325 |
	file -
