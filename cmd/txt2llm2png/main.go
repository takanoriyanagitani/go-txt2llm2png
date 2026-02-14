package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	oa "github.com/ollama/ollama/api"
	tp "github.com/takanoriyanagitani/go-txt2llm2png"
)

var (
	ErrModelMissing  error = errors.New("model is required")
	ErrPromptMissing error = errors.New("prompt is required (either via --prompt flag or stdin)")
)

func sub(ctx context.Context, args []string) error {
	var model string
	var prompt string
	var width int
	var height int
	var steps int
	var seed int

	flag.StringVar(&model, "model", "", "Model name (e.g., x/flux2-klein:4b-bf16)")
	flag.StringVar(&prompt, "prompt", "", "Prompt for image generation (e.g., 'draw a dog')")
	flag.IntVar(&width, "width", int(tp.RequestDefault.Width), "Image width")
	flag.IntVar(&height, "height", int(tp.RequestDefault.Height), "Image height")
	flag.IntVar(&steps, "steps", int(tp.RequestDefault.Steps), "Number of steps")
	flag.IntVar(&seed, "seed", tp.RequestDefault.Seed, "Seed for random generation")

	// Parse flags from the provided args, excluding the program name
	err := flag.CommandLine.Parse(args[1:])
	if nil != err {
		flag.Usage()
		return nil //nolint:nilerr
	}

	if model == "" {
		return ErrModelMissing
	}

	if prompt == "" {
		// If prompt is empty, try to read from stdin
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read prompt from stdin: %w", err)
		}
		prompt = strings.TrimSpace(buf.String())
	}

	if prompt == "" {
		return ErrPromptMissing
	}

	ocli, e := oa.ClientFromEnvironment()
	if nil != e {
		return e
	}

	cli := tp.Client{Client: ocli}
	req := tp.RequestDefault.
		WithModel(model).
		WithPrompt(prompt).
		WithWidth(int32(width)).
		WithHeight(int32(height)).
		WithSteps(int32(steps)).
		WithSeed(seed)
	generated, err := cli.Generate(ctx, req)
	if nil != err {
		return err
	}

	return generated.Write(os.Stdout)
}

func main() {
	e := sub(context.Background(), os.Args)
	if nil != e {
		log.Printf("%v\n", e)
		os.Exit(1) // Exit with a non-zero status code on error
	}
}
