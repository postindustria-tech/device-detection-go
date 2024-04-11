package pipeline

import "github.com/51Degrees/device-detection-go/v4/dd"

type Pipeline struct {
	logger                        logWrapper
	dataFile                      string
	licenceKey                    string
	createTempDataCopy            bool
	dataUpdateOnStartEnabled      bool
	dataFileUrl                   string
	dataFilePullEveryMs           int
	isScheduledFilePullingEnabled bool
	loggerEnabled                 bool
	manager                       *dd.ResourceManager
	config                        *dd.ConfigHash
}

// NewResultsHash creates a new ResultsHash object
// uaCapacity is the initial capacity of the User-Agent cache
// overridesCapacity is the initial capacity of the overrides cache
func (p *Pipeline) NewResultsHash(uaCapacity uint32, overridesCapacity uint32) *dd.ResultsHash {
	return dd.NewResultsHash(p.manager, uaCapacity, overridesCapacity)
}

// Run starts the pipeline
func (p *Pipeline) Run() error {
	dataFile, err := dd.GetFilePath(".", []string{p.dataFile})
	if err != nil {
		return err
	}
	p.dataFile = dataFile

	// TODO: what are properties of dd.ResourceManager?
	err = dd.InitManagerFromFile(p.manager, *p.config, "", p.dataFile)
	if err != nil {
		return err
	}

	if p.isScheduledFilePullingEnabled {
		go p.scheduleFilePulling()
	}

	return nil
}

type PipelineOptions func(cfg *Pipeline) error

// what should this do? c# has it and it's used in the example.
func UseOnPremise(dataFile, licenceKey string, createTempDataCopy bool) PipelineOptions {

	return func(cfg *Pipeline) error {
		cfg.dataFile = dataFile
		cfg.licenceKey = licenceKey
		cfg.createTempDataCopy = createTempDataCopy

		return nil
	}
}

// WithDataUpdateUrl sets the URL to pull data from and the interval in milliseconds
func WithDataUpdateUrl(url string, everyMs int) PipelineOptions {
	return func(cfg *Pipeline) error {
		cfg.dataFileUrl = url
		cfg.dataFilePullEveryMs = everyMs
		cfg.isScheduledFilePullingEnabled = true

		return nil
	}
}

// ToggleLogger enables or disables the logger
func ToggleLogger(enabled bool) PipelineOptions {
	return func(cfg *Pipeline) error {
		cfg.logger.enabled = enabled
		return nil
	}
}

// WithCustomLogger sets a custom logger
func WithCustomLogger(logger LogWriter) PipelineOptions {
	return func(cfg *Pipeline) error {
		cfg.logger = logWrapper{
			enabled: true,
			logger:  logger,
		}

		return nil
	}

}

func New(manager *dd.ResourceManager, config *dd.ConfigHash, opts ...PipelineOptions) (*Pipeline, error) {
	pl := &Pipeline{
		logger: logWrapper{
			logger:  DefaultLogger,
			enabled: true,
		},
		manager: manager,
		config:  config,
	}

	for _, opt := range opts {
		err := opt(pl)
		if err != nil {
			return nil, err
		}
	}

	return pl, nil
}
