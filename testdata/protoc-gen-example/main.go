package main

import (
	pgs "github.com/joakim-brannstrom/protoc-gen-star"
	pgsgo "github.com/joakim-brannstrom/protoc-gen-star/lang/go"
)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		ASTPrinter(),
		JSONify(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}
