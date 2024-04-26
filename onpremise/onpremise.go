package onpremise

import (
	"errors"
	"fmt"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"net/url"
	"os"
	"path/filepath"
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
	rdySignal                     chan error
	stopCh                        chan struct{}
	fileSynced                    bool
	isManagerInitialized          bool
	product                       string
	maxRetries                    int
}

const (
	defaultDataFileUrl = "https://distributor.51degrees.com/api/v2/download?Type=HashV41&Download=True&Product=V4TAC"
)

// run starts the engine
func (p *Engine) run() error {
	if len(p.dataFile) > 0 {
		err := p.initializeManager()
		if err != nil {
			return err
		}
	}

	err := p.validateAndAppendUrlParams()
	if err != nil {
		return err
	}

	if p.isScheduledFilePullingEnabled {
		go p.scheduleFilePulling()
	}

	return nil
}

func (p *Engine) initializeManager() error {

	err := dd.InitManagerFromFile(p.manager, *p.config, "", p.dataFile)
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
			p.dataFile = filepath.FromSlash("data.hash")
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
		path := filepath.Join(path)
		_, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to get file path: %w", err)
		}

		cfg.dataFile = path
		return nil
	}
}

func SetLicenceKey(key string) EngineOptions {
	return func(cfg *Engine) error {
		if !cfg.isDefaultDataFileUrl() {
			return errors.New("licence key can only be set when using default data file url")
		}
		cfg.licenceKey = key
		cfg.isScheduledFilePullingEnabled = true
		return nil
	}
}

func SetProduct(product string) EngineOptions {
	return func(cfg *Engine) error {
		if !cfg.isDefaultDataFileUrl() {
			return errors.New("product can only be set when using default data file url")
		}

		cfg.product = product
		cfg.isScheduledFilePullingEnabled = true
		return nil
	}
}

// WithDataUpdateUrl sets the URL to pull data from and the interval in milliseconds
func WithDataUpdateUrl(urlStr string, seconds int) EngineOptions {
	return func(cfg *Engine) error {
		_, err := url.ParseRequestURI(urlStr)
		if err != nil {
			return err
		}

		cfg.dataFileUrl = urlStr
		cfg.dataFilePullEveryMs = seconds * 1000
		cfg.isScheduledFilePullingEnabled = true

		return nil
	}
}

func SetMaxRetries(retries int) EngineOptions {
	return func(cfg *Engine) error {
		cfg.maxRetries = retries
		return nil
	}
}

func SetPollingInterval(seconds int) EngineOptions {
	return func(cfg *Engine) error {
		cfg.dataFilePullEveryMs = seconds * 1000
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
		manager:     manager,
		config:      config,
		rdySignal:   make(chan error, 1),
		stopCh:      make(chan struct{}, 1),
		fileSynced:  false,
		dataFileUrl: defaultDataFileUrl,
		//default 15 minutes
		dataFilePullEveryMs: 15 * 60 * 1000,
	}

	for _, opt := range opts {
		err := opt(pl)
		if err != nil {
			pl.Stop()
			return nil, err
		}
	}

	err := pl.run()
	if err != nil {
		pl.Stop()
		return nil, err
	}

	if pl.isScheduledFilePullingEnabled && !pl.isManagerInitialized {
		//wait for first file pull
		err := <-pl.rdySignal
		defer close(pl.rdySignal)
		if err != nil {
			pl.Stop()
			return nil, err
		}
	}

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
	defer evidenceHash.Free()

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
	p.stopCh <- struct{}{}
	if p.manager != nil {
		p.manager.Free()
	}
	close(p.stopCh)
}

func (p *Engine) GetHttpHeaderKeys() []dd.EvidenceKey {
	return p.manager.HttpHeaderKeys
}

func (p *Engine) appendLicenceKey() error {
	urlParsed, err := url.Parse(p.dataFileUrl)
	if err != nil {
		return err
	}
	query := urlParsed.Query()
	query.Set("LicenseKeys", p.licenceKey)
	urlParsed.RawQuery = query.Encode()

	p.dataFileUrl = urlParsed.String()

	return nil
}

func (p *Engine) isDefaultDataFileUrl() bool {
	return p.dataFileUrl == defaultDataFileUrl
}

func (p *Engine) appendProduct() error {
	urlParsed, err := url.Parse(p.dataFileUrl)
	if err != nil {
		return fmt.Errorf("failed to parse data file url: %w", err)
	}
	query := urlParsed.Query()
	query.Set("Product", p.product)
	urlParsed.RawQuery = query.Encode()

	p.dataFileUrl = urlParsed.String()

	return nil
}

var (
	ErrLicenceKeyAndProductRequired = errors.New("licence key and product are required")
)

func (p *Engine) validateAndAppendUrlParams() error {
	if p.isDefaultDataFileUrl() && !p.hasDefaultDistributorParams() && p.isScheduledFilePullingEnabled {
		return ErrLicenceKeyAndProductRequired
	} else if p.isDefaultDataFileUrl() && p.isScheduledFilePullingEnabled {
		err := p.appendLicenceKey()
		if err != nil {
			return err
		}
		err = p.appendProduct()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Engine) hasDefaultDistributorParams() bool {
	return len(p.licenceKey) > 0 && len(p.product) > 0
}

func (p *Engine) hasSomeDistributorParams() bool {
	return len(p.licenceKey) > 0 || len(p.product) > 0
}
