#!/bin/sh

bname="txt2llm2png"
bdir="./cmd/${bname}"
oname="${bdir}/${bname}"

mkdir -p "${bdir}"

go \
	build \
	-v \
	./...

go \
	build \
	-v \
	-o "${oname}" \
	"${bdir}"
