package registry

//
// Registers the `nym' experiment.
//

import (
	"github.com/ooni/probe-cli/v3/internal/experiment/nym"
	"github.com/ooni/probe-cli/v3/internal/model"
)

func init() {
	AllExperiments["nym"] = &Factory{
		build: func(config interface{}) model.ExperimentMeasurer {
			return nym.NewExperimentMeasurer(
				*config.(*nym.Config),
			)
		},
		config:      &nym.Config{},
		inputPolicy: model.InputNone,
	}
}
