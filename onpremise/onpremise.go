package onpremise

import (
	"fmt"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"os"
)

type Engine struct {
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
	totalFilePulls                int
	rdySignal                     chan struct{}
	fileSynced                    bool
	isManagerInitialized          bool
}

// run starts the engine
func (p *Engine) run() error {
	if len(p.dataFile) > 0 {
		err := p.initializeManager()
		if err != nil {
			return err
		}
	}

	if p.isScheduledFilePullingEnabled {
		go p.scheduleFilePulling()
	}

	return nil
}

func (p *Engine) initializeManager() error {
	dataFile, err := dd.GetFilePath(".", []string{p.dataFile})
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}
	p.dataFile = dataFile

	err = dd.InitManagerFromFile(p.manager, *p.config, "", p.dataFile)
	if err != nil {
		return fmt.Errorf("failed to init manager from file: %w", err)
	}
	p.isManagerInitialized = true

	return nil
}

func (p *Engine) createDatafileIfNotExists() error {
	_, err := os.Stat(p.dataFile)
	if err != nil {
		if len(p.dataFile) == 0 {
			p.dataFile = "data.hash"
		}

		_, err = os.Create(p.dataFile)
		if err != nil {
			return fmt.Errorf("failed to create data file: %w", err)
		}
	}

	return nil
}

type EngineOptions func(cfg *Engine) error

func WithDataFile(path string) EngineOptions {
	return func(cfg *Engine) error {
		cfg.dataFile = path
		return nil
	}
}

// WithDataUpdateUrl sets the URL to pull data from and the interval in milliseconds
func WithDataUpdateUrl(url string, everyMs int) EngineOptions {
	return func(cfg *Engine) error {
		cfg.dataFileUrl = url
		cfg.dataFilePullEveryMs = everyMs
		cfg.isScheduledFilePullingEnabled = true

		return nil
	}
}

// ToggleLogger enables or disables the logger
func ToggleLogger(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.logger.enabled = enabled
		return nil
	}
}

// WithCustomLogger sets a custom logger
func WithCustomLogger(logger LogWriter) EngineOptions {
	return func(cfg *Engine) error {
		cfg.logger = logWrapper{
			enabled: true,
			logger:  logger,
		}

		return nil
	}

}

func New(config *dd.ConfigHash, opts ...EngineOptions) (*Engine, error) {
	manager := dd.NewResourceManager()

	pl := &Engine{
		logger: logWrapper{
			logger:  DefaultLogger,
			enabled: true,
		},
		manager:    manager,
		config:     config,
		rdySignal:  make(chan struct{}, 1),
		fileSynced: false,
	}

	for _, opt := range opts {
		err := opt(pl)
		if err != nil {
			return nil, err
		}
	}

	err := pl.run()
	if err != nil {
		return nil, err
	}

	//wait for first file pull
	<-pl.rdySignal
	close(pl.rdySignal)

	return pl, nil
}

type Evidence struct {
	Prefix dd.EvidencePrefix
	Key    string
	Value  string
}

func (p *Engine) Process(evidences []Evidence) (*dd.ResultsHash, error) {
	evidenceHash, err := mapEvidence(evidences)
	if err != nil {
		return nil, err
	}

	results := dd.NewResultsHash(p.manager, uint32(evidenceHash.Count()), 0)
	err = results.MatchEvidence(evidenceHash)
	if err != nil {
		return nil, fmt.Errorf("failed to match evidence: %w", err)
	}

	return results, nil
}

func mapEvidence(evidences []Evidence) (*dd.Evidence, error) {
	evidenceHash := dd.NewEvidenceHash(uint32(len(evidences)))

	for _, evidence := range evidences {
		err := evidenceHash.Add(evidence.Prefix, evidence.Key, evidence.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to add evidence: %w", err)
		}
	}

	return evidenceHash, nil
}

func (p *Engine) Stop() {
	p.manager.Free()
}
