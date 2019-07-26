package cli

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"

	v2 "asyncapi-converter/pkg/converter/v2"
	"asyncapi-converter/pkg/decode"
	asyncapiEncode "asyncapi-converter/pkg/encode"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	errFileDoesNotExist = errors.New("file does not exist")
	errInvalidArgument  = errors.New("invalid argument")
)

type Encode = func(interface{}, io.Writer) error

type Converter interface {
	Convert(reader io.Reader, writer io.Writer) error
}

type cli struct {
	docopt.Opts
	data interface{}
}

func New(opts docopt.Opts) cli {
	return cli{
		Opts: opts,
	}
}

func (h cli) id() *string {
	idOption, ok := h.Opts["--id"]
	if !ok || idOption == nil {
		return nil
	}
	id := fmt.Sprintf("%v", idOption)
	return &id
}

func (h cli) encode() (Encode, error) {
	if _, ok := h.Opts["--toYAML"]; !ok {
		return asyncapiEncode.JsonEncoder, nil
	}
	toYaml, ok := h.Opts["--toYAML"].(bool)
	if !ok {
		return nil, errors.Wrap(errInvalidArgument, "--toYAML")
	}
	if toYaml {
		return asyncapiEncode.YamlEncoder, nil
	}
	return asyncapiEncode.JsonEncoder, nil
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (h cli) reader() (io.Reader, error) {
	fileOption := h.Opts["<PATH>"]
	path := fmt.Sprintf("%v", fileOption)
	if isUrl(path) {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Wrap(errFileDoesNotExist, path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (h cli) NewConverterAndReader() (Converter, io.Reader, error) {
	reader, err := h.reader()
	if err != nil {
		return nil, nil, err
	}
	encode, err := h.encode()
	if err != nil {
		return nil, nil, err
	}
	converter, err := v2.NewConverter(decode.JsonDecoderWithYamlFallback, encode, v2.WithId(h.id()))
	return converter, reader, err
}
