package txt2llm2png

import (
	"context"
	"encoding/base64"
	"io"
	"strings"

	oa "github.com/ollama/ollama/api"
)

type Generated struct{ oa.GenerateResponse }

func (g Generated) ImageBase64() string {
	return g.GenerateResponse.Image
}

func (g Generated) WriteBase64(wtr io.Writer) error {
	var img string = g.ImageBase64()
	_, e := io.WriteString(wtr, img)
	return e
}

func (g Generated) Write(wtr io.Writer) error {
	var rdr io.Reader = strings.NewReader(g.GenerateResponse.Image)
	var dec io.Reader = base64.NewDecoder(
		base64.StdEncoding,
		rdr,
	)
	_, e := io.Copy(wtr, dec)
	return e
}

type Request struct {
	Model    string
	Prompt   string
	System   string
	Template string

	Width  int32
	Height int32
	Steps  int32

	Seed int
}

//nolint:gochecknoglobals
var RequestDefault Request = Request{
	Width:  128,
	Height: 128,
	Steps:  4,
	Seed:   -1,
}

func (r Request) ToOptionsMap() map[string]any {
	return map[string]any{
		"seed": r.Seed,
	}
}

func (r Request) ToGenRequest() oa.GenerateRequest {
	return oa.GenerateRequest{
		Model:    r.Model,
		Prompt:   r.Prompt,
		System:   r.System,
		Template: r.Template,
		Stream:   new(bool),
		Options:  r.ToOptionsMap(),
		Width:    r.Width,
		Height:   r.Height,
		Steps:    r.Steps,
	}
}

func (r Request) WithPrompt(txt string) Request {
	r.Prompt = txt
	return r
}

func (r Request) WithModel(model string) Request {
	r.Model = model
	return r
}

func (r Request) WithHeight(h int32) Request {
	r.Height = h
	return r
}

func (r Request) WithWidth(w int32) Request {
	r.Width = w
	return r
}

func (r Request) WithSteps(s int32) Request {
	r.Steps = s
	return r
}

func (r Request) WithSeed(s int) Request {
	r.Seed = s
	return r
}

type Client struct{ *oa.Client }

func (c Client) Generate(ctx context.Context, req Request) (Generated, error) {
	var greq oa.GenerateRequest = req.ToGenRequest()
	var gres oa.GenerateResponse
	err := c.Client.Generate(
		ctx,
		&greq,
		func(res oa.GenerateResponse) error {
			gres = res
			return nil
		},
	)
	return Generated{gres}, err
}
