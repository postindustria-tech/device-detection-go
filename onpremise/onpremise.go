package onpremise

import (
	"errors"
	"fmt"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/google/uuid"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Engine struct {
	sync.RWMutex
	logger                        logWrapper
	fileWatcher                   fileWatcher
	dataFile                      string
	licenceKey                    string
	dataFileUrl                   string
	dataFilePullEveryMs           int
	isScheduledFilePullingEnabled bool
	isAutoUpdateEnabled           bool
	loggerEnabled                 bool
	manager                       *dd.ResourceManager
	config                        *dd.ConfigHash
	totalFilePulls                int
	stopCh                        chan struct{}
	fileSynced                    bool
	product                       string
	maxRetries                    int
	lastModificationTimestamp     *time.Time
	isFileWatcherEnabled          bool
	isUpdateOnStartEnabled        bool
	isCreateTempDataCopyEnabled   bool
	tempDataFile                  string
	tempDataDir                   string
	isCopyingFile                 bool
	randomization                 int
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
	if e.isCreateTempDataCopyEnabled {
		dirpath, tempfilepath, err := e.copyToTempFile()
		if err != nil {
			return err
		}
		tempDataPathFull := filepath.Join(dirpath, tempfilepath)
		e.manager = dd.NewResourceManager()
		err = dd.InitManagerFromFile(e.manager, *e.config, "", tempDataPathFull)
		if err != nil {
			return err
		}
		e.tempDataFile = tempfilepath
	} else {
		e.manager = dd.NewResourceManager()
		err := dd.InitManagerFromFile(e.manager, *e.config, "", e.dataFile)
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

type EngineOptions func(cfg *Engine) error

// WithDataFile sets the path to the local data file, this parameter is required to start the engine
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

// SetLicenceKey sets the licence key to use when pulling the data file
// this option can only be used when using the default data file url from 51Degrees, it will be appended as a query parameter
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

// SetProduct sets the product to use when pulling the data file
// this option can only be used when using the default data file url from 51Degrees, it will be appended as a query parameter
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

// SetMaxRetries sets the maximum number of retries to pull the data file
func SetMaxRetries(retries int) EngineOptions {
	return func(cfg *Engine) error {
		cfg.maxRetries = retries
		return nil
	}
}

// SetPollingInterval sets the interval in seconds to pull the data file
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

// ToggleUpdateOnStart enables or disables update on start
// default is true
// if enabled, engine will pull the data file from the distributor on start of the engine
func ToggleUpdateOnStart(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.isUpdateOnStartEnabled = enabled

		return nil
	}
}

// ToggleAutoUpdate enables or disables auto update
// default is true
// if enabled, engine will automatically pull the data file from the distributor
// if disabled, engine will not pull the data file from the distributor
// options like WithDataUpdateUrl, SetLicenceKey will be ignored since auto update is disabled
func ToggleAutoUpdate(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.isAutoUpdateEnabled = enabled

		return nil
	}
}

// ToggleCreateTempDataCopy enables or disables creating a temp copy of the data file
// default is true
// if enabled, engine will create a temp copy of the data file and use it to initialize the manager
// if disabled, engine will use the original data file to initialize the manager
// this is useful when 3rd party updates the data file on file system
func ToggleCreateTempDataCopy(enabled bool) EngineOptions {
	return func(cfg *Engine) error {
		cfg.isCreateTempDataCopyEnabled = enabled

		return nil
	}
}

// SetTempDataDir sets the directory to store the temp data file
// default is current directory
func SetTempDataDir(dir string) EngineOptions {
	return func(cfg *Engine) error {
		dirFileInfo, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("failed to get file path: %w", err)
		}

		if !dirFileInfo.IsDir() {
			return fmt.Errorf("path is not a directory: %s", dir)
		}

		cfg.tempDataDir = dir
		return nil
	}
}

// Randomization sets the randomization time in seconds
// default is 0
// if set, when scheduling the file pulling, it will add randomization time to the interval
// this is useful to avoid all engines pulling the data file at the same time in case of multiple engines/instances
func Randomization(seconds int) EngineOptions {
	return func(cfg *Engine) error {
		cfg.randomization = seconds * 1000
		return nil
	}
}

func New(config *dd.ConfigHash, opts ...EngineOptions) (*Engine, error) {
	pl := &Engine{
		logger: logWrapper{
			logger:  DefaultLogger,
			enabled: true,
		},
		config:      config,
		stopCh:      make(chan struct{}, 1),
		fileSynced:  false,
		dataFileUrl: defaultDataFileUrl,
		//default 15 minutes
		dataFilePullEveryMs:         30 * 60 * 1000,
		isFileWatcherEnabled:        true,
		isUpdateOnStartEnabled:      true,
		isAutoUpdateEnabled:         true,
		isCreateTempDataCopyEnabled: true,
		tempDataDir:                 ".",
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
	e.RLock()
	defer e.RUnlock()
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
	if e.isCreateTempDataCopyEnabled {
		os.Remove(e.tempDataFile)
	}
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
	e.Lock()
	defer e.Unlock()
	if e.isCreateTempDataCopyEnabled {
		dirPath, filePath, err := e.copyToTempFile()
		if err != nil {
			e.logger.Printf("failed to copy data file: %v", err)
		}
		fullPath := filepath.Join(dirPath, filePath)
		err = e.replaceManager(fullPath)
		if err != nil {
			e.logger.Printf("failed to replace manager: %v", err)
		}

		oldFullPath := filepath.Join(e.tempDataDir, e.tempDataFile)
		err = os.Remove(oldFullPath)
		if err != nil {
			e.logger.Printf("failed to remove old temp data file: %v", err)
		}
	} else {
		err := e.replaceManager(e.dataFile)
		if err != nil {
			e.logger.Printf("failed to replace manager: %v", err)
		}
	}
}

func (e *Engine) replaceManager(filePath string) error {
	if !e.isCreateTempDataCopyEnabled {
		err := e.manager.ReloadFromOriginalFile()
		if err != nil {
			return fmt.Errorf("failed to reload manager from original file: %w", err)
		}
		return nil
	}

	newManager := dd.NewResourceManager()
	err := dd.InitManagerFromFile(newManager, *e.config, "", filePath)
	if err != nil {
		return fmt.Errorf("failed to init manager from file: %w", err)

	}
	oldManager := e.manager
	oldManager.Free()
	e.manager = newManager

	return nil
}

func newTempFilePath(originalName string) string {
	randomUuid := uuid.NewString()
	newName := fmt.Sprintf("%s-%s", randomUuid, originalName)

	return newName
}

func (e *Engine) copyToTempFile() (string, string, error) {
	data, err := os.ReadFile(e.dataFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to read data file: %w", err)
	}
	_, originalFileName := filepath.Split(e.dataFile)

	tempFileName := newTempFilePath(originalFileName)

	path := filepath.Join(e.tempDataDir, tempFileName)

	f, err := os.Create(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp data file: %w", err)
	}
	defer f.Close()

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write temp data file: %w", err)
	}

	return e.tempDataDir, tempFileName, nil
}

func (e *Engine) getFilePath() string {
	if e.isCreateTempDataCopyEnabled {
		return filepath.Join(e.tempDataDir, e.tempDataFile)
	}

	return e.dataFile
}
