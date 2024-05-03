package onpremise

import (
	"errors"
	"fmt"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Engine struct {
	logger                        logWrapper
	fileWatcher                   fileWatcher
	dataFile                      string
	licenceKey                    string
	createTempDataCopy            bool
	dataFileUrl                   string
	dataFilePullEveryMs           int
	isScheduledFilePullingEnabled bool
	isAutoUpdateEnabled           bool
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
	lastModificationTimestamp     *time.Time
	isFileWatcherEnabled          bool
	isUpdateOnStartEnabled        bool
}

const (
	defaultDataFileUrl = "https://distributor.51degrees.com/api/v2/download?Type=HashV41&Download=True&Product=V4TAC"
)

var (
	ErrNoDataFileProvided           = errors.New("no data file provided")
	ErrTooManyRetries               = errors.New("too many retries to pull data file")
	ErrFileNotModified              = errors.New("data file not modified")
	ErrLicenceKeyAndProductRequired = errors.New("licence key and product are required")
)

// run starts the engine
func (e *Engine) run() error {
	if len(e.dataFile) > 0 {
		err := e.initializeManager()
		if err != nil {
			return err
		}
	}

	err := e.validateAndAppendUrlParams()
	if err != nil {
		return err
	}

	if e.isScheduledFilePullingEnabled && e.isAutoUpdateEnabled {
		go e.scheduleFilePulling()
	}

	return nil
}

func (e *Engine) initializeManager() error {

	err := dd.InitManagerFromFile(e.manager, *e.config, "", e.dataFile)
	if err != nil {
		return fmt.Errorf("failed to init manager from file: %w", err)
	}
	e.isManagerInitialized = true

	return nil
}

func (e *Engine) createDatafileIfNotExists() error {
	_, err := os.Stat(e.dataFile)
	if err != nil {
		if len(e.dataFile) == 0 {
			e.dataFile, err = filepath.Abs("data.hash")
			if err != nil {
				e.dataFile = filepath.FromSlash("data.hash")
			}
		}

		_, err = os.Create(e.dataFile)
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

// ToggleFileWatch enables or disables file watching
// in case 3rd party updates the data file on file system
// engine will automatically reload the data file
// default is true
func ToggleFileWatch(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.isFileWatcherEnabled = enabled
		return nil
	}
}

func ToggleUpdateOnStart(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.isUpdateOnStartEnabled = enabled

		return nil
	}
}

func ToggleAutoUpdate(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.isAutoUpdateEnabled = enabled

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
		dataFilePullEveryMs:    30 * 60 * 1000,
		isFileWatcherEnabled:   true,
		isUpdateOnStartEnabled: true,
		isAutoUpdateEnabled:    true,
	}

	for _, opt := range opts {
		err := opt(pl)
		if err != nil {
			pl.Stop()
			return nil, err
		}
	}

	if pl.dataFile == "" {
		return nil, ErrNoDataFileProvided
	}

	err := pl.run()
	if err != nil {
		pl.Stop()
		return nil, err
	}

	if pl.isScheduledFilePullingEnabled && !pl.isManagerInitialized && pl.isAutoUpdateEnabled {
		//wait for first file pull
		err := <-pl.rdySignal
		defer close(pl.rdySignal)
		if err != nil {
			pl.Stop()
			return nil, err
		}
	}

	// if file watcher is enabled, start the watcher
	if pl.isFileWatcherEnabled {
		pl.fileWatcher, err = newFileWatcher(pl.logger)
		if err != nil {
			return nil, err
		}
		// this will watch the data file, if it changes, it will reload the data file in the manager
		err = pl.fileWatcher.watch(pl.dataFile, pl.handleFileExternallyChanged)
		if err != nil {
			return nil, err
		}
		go pl.fileWatcher.run()
	}

	return pl, nil
}

type Evidence struct {
	Prefix dd.EvidencePrefix
	Key    string
	Value  string
}

func (e *Engine) Process(evidences []Evidence) (*dd.ResultsHash, error) {
	evidenceHash, err := mapEvidence(evidences)
	if err != nil {
		return nil, err
	}
	defer evidenceHash.Free()

	results := dd.NewResultsHash(e.manager, uint32(evidenceHash.Count()), 0)
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

func (e *Engine) Stop() {
	e.stopCh <- struct{}{}
	if e.manager != nil {
		e.manager.Free()
	} else {
		e.logger.Printf("manager is nil")
	}
	close(e.stopCh)
}

func (e *Engine) GetHttpHeaderKeys() []dd.EvidenceKey {
	return e.manager.HttpHeaderKeys
}

func (e *Engine) appendLicenceKey() error {
	urlParsed, err := url.Parse(e.dataFileUrl)
	if err != nil {
		return err
	}
	query := urlParsed.Query()
	query.Set("LicenseKeys", e.licenceKey)
	urlParsed.RawQuery = query.Encode()

	e.dataFileUrl = urlParsed.String()

	return nil
}

func (e *Engine) isDefaultDataFileUrl() bool {
	return e.dataFileUrl == defaultDataFileUrl
}

func (e *Engine) appendProduct() error {
	urlParsed, err := url.Parse(e.dataFileUrl)
	if err != nil {
		return fmt.Errorf("failed to parse data file url: %w", err)
	}
	query := urlParsed.Query()
	query.Set("Product", e.product)
	urlParsed.RawQuery = query.Encode()

	e.dataFileUrl = urlParsed.String()

	return nil
}

func (e *Engine) validateAndAppendUrlParams() error {
	if e.isDefaultDataFileUrl() && !e.hasDefaultDistributorParams() && e.isScheduledFilePullingEnabled {
		return ErrLicenceKeyAndProductRequired
	} else if e.isDefaultDataFileUrl() && e.isScheduledFilePullingEnabled {
		err := e.appendLicenceKey()
		if err != nil {
			return err
		}
		err = e.appendProduct()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) hasDefaultDistributorParams() bool {
	return len(e.licenceKey) > 0 && len(e.product) > 0
}

func (e *Engine) hasSomeDistributorParams() bool {
	return len(e.licenceKey) > 0 || len(e.product) > 0
}

func (e *Engine) handleFileExternallyChanged() {
	err := e.manager.ReloadFromOriginalFile()
	if err != nil {
		e.logger.Printf("failed to reload data file: %v", err)
	}
}
