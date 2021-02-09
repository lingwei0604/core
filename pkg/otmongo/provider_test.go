package otmongo

import (
	"github.com/DoNewsCode/std/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"testing"
)

func TestNewMongoFactory(t *testing.T) {
	factory, cleanup := ProvideMongoFactory(MongoParam{
		In: dig.In{},
		Conf: config.MapAdapter{"mongo": map[string]struct{ Uri string }{
			"default": {
				Uri: "mongodb://127.0.0.1:27017",
			},
			"alternative": {
				Uri: "mongodb://127.0.0.1:27017",
			},
		}},
		Tracer: nil,
	})
	def, err := factory.Make("default")
	assert.NoError(t, err)
	assert.NotNil(t, def)
	alt, err := factory.Make("alternative")
	assert.NoError(t, err)
	assert.NotNil(t, alt)
	assert.NotNil(t, cleanup)
	cleanup()
}