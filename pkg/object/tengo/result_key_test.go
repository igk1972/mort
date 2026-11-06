package tengo

import (
	"net/url"
	"testing"

	"github.com/aldor007/mort/pkg/config"
	"github.com/aldor007/mort/pkg/object"
	"github.com/stretchr/testify/assert"
)

func resultKeyConfig(t *testing.T, resultKey string) *config.Config {
	mortConfig := config.Config{BaseConfigPath: "./testdata"}
	err := mortConfig.LoadFromString(`
buckets:
    convert:
        transform:
            kind: "tengo"
            tengoPath: "result-key.tengo"
            parentBucket: "originals"
            resultKey: "` + resultKey + `"
        storages:
            basic:
                kind: "noop"
            transform:
                kind: "noop"
    originals:
        storages:
            basic:
                kind: "noop"
`)
	assert.Nil(t, err)
	return &mortConfig
}

func parseURL(t *testing.T, mortConfig *config.Config, raw string) (*object.FileObject, error) {
	u, err := url.Parse(raw)
	assert.Nil(t, err)
	return object.NewFileObject(u, mortConfig)
}

func TestScriptResultKeyStoresUnderScriptKey(t *testing.T) {
	obj, err := parseURL(t, resultKeyConfig(t, "script"), "http://mort/convert/aa/bb/id.png")

	assert.Nil(t, err)
	assert.Equal(t, "/aa/bb/id.jpg", obj.Key)
	assert.Equal(t, "/aa/bb/id.png", obj.Parent.Key)
}

func TestScriptResultKeyIgnoredWithoutScriptMode(t *testing.T) {
	obj, err := parseURL(t, resultKeyConfig(t, ""), "http://mort/convert/aa/bb/id.png")

	assert.Nil(t, err)
	assert.Equal(t, "/aa/bb/id.png", obj.Key)
}

func TestScriptResultKeyRequired(t *testing.T) {
	_, err := parseURL(t, resultKeyConfig(t, "script"), "http://mort/convert/aa/bb/id.png?mode=missing")

	assert.ErrorContains(t, err, "did not set obj.resultKey")
}

func TestScriptResultKeyRejectsParentSegments(t *testing.T) {
	_, err := parseURL(t, resultKeyConfig(t, "script"), "http://mort/convert/aa/bb/id.png?mode=escape")

	assert.ErrorContains(t, err, "must not contain '..'")
}

func TestScriptResultKeyMustBeString(t *testing.T) {
	_, err := parseURL(t, resultKeyConfig(t, "script"), "http://mort/convert/aa/bb/id.png?mode=number")

	assert.ErrorContains(t, err, "index value type: int")
}

func TestScriptResultKeyRequiresTengoKind(t *testing.T) {
	mortConfig := config.Config{}
	err := mortConfig.LoadFromString(`
buckets:
    convert:
        transform:
            kind: "query"
            resultKey: "script"
        storages:
            basic:
                kind: "noop"
            transform:
                kind: "noop"
`)

	assert.ErrorContains(t, err, `resultKey "script" requires kind tengo`)
}
