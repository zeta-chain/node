package debug

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	stderrors "github.com/pkg/errors"

	"github.com/cosmos/evm/rpc/backend"
	rpctypes "github.com/cosmos/evm/rpc/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/server"
)

// HandlerT keeps track of the cpu profiler and trace execution
type HandlerT struct {
	cpuFilename   string
	cpuFile       io.WriteCloser
	mu            sync.Mutex
	traceFilename string
	traceFile     io.WriteCloser
}

// API is the collection of tracing APIs exposed over the private debugging endpoint.
type API struct {
	ctx              *server.Context
	logger           log.Logger
	profilingEnabled bool
	backend          backend.EVMBackend
	handler          *HandlerT
}

// NewAPI creates a new API definition for the tracing methods of the Ethereum service.
func NewAPI(
	ctx *server.Context,
	backend backend.EVMBackend,
	profilingEnabled bool,
) *API {
	return &API{
		ctx:              ctx,
		logger:           ctx.Logger.With("module", "debug"),
		backend:          backend,
		profilingEnabled: profilingEnabled,
		handler:          new(HandlerT),
	}
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (a *API) TraceTransaction(hash common.Hash, config *rpctypes.TraceConfig) (interface{}, error) {
	a.logger.Debug("debug_traceTransaction", "hash", hash)
	return a.backend.TraceTransaction(hash, config)
}

// TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (a *API) TraceBlockByNumber(height rpctypes.BlockNumber, config *rpctypes.TraceConfig) ([]*evmtypes.TxTraceResult, error) {
	a.logger.Debug("debug_traceBlockByNumber", "height", height)
	if height == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	// Get CometBFT Block
	resBlock, err := a.backend.CometBlockByNumber(height)
	if err != nil {
		a.logger.Debug("get block failed", "height", height, "error", err.Error())
		return nil, err
	}

	return a.backend.TraceBlock(rpctypes.BlockNumber(resBlock.Block.Height), config, resBlock)
}

// TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (a *API) TraceBlockByHash(hash common.Hash, config *rpctypes.TraceConfig) ([]*evmtypes.TxTraceResult, error) {
	a.logger.Debug("debug_traceBlockByHash", "hash", hash)
	// Get CometBFT Block
	resBlock, err := a.backend.CometBlockByHash(hash)
	if err != nil {
		a.logger.Debug("get block failed", "hash", hash.Hex(), "error", err.Error())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		a.logger.Debug("block not found", "hash", hash.Hex())
		return nil, errors.New("block not found")
	}

	return a.backend.TraceBlock(rpctypes.BlockNumber(resBlock.Block.Height), config, resBlock)
}

// TraceBlock returns the structured logs created during the execution of
// EVM and returns them as a JSON object. It accepts an RLP-encoded block.
func (a *API) TraceBlock(tblockRlp hexutil.Bytes, config *rpctypes.TraceConfig) ([]*evmtypes.TxTraceResult, error) {
	a.logger.Debug("debug_traceBlock", "size", len(tblockRlp))
	// Decode RLP-encoded block
	var block types.Block
	if err := rlp.DecodeBytes(tblockRlp, &block); err != nil {
		a.logger.Debug("failed to decode block", "error", err.Error())
		return nil, fmt.Errorf("could not decode block: %w", err)
	}

	// Get block number from the decoded block
	blockNum := block.NumberU64()
	if blockNum > math.MaxInt64 {
		return nil, fmt.Errorf("block number overflow: %d exceeds max int64", blockNum)
	}
	blockNumber := rpctypes.BlockNumber(blockNum) //#nosec G115 -- overflow checked above
	a.logger.Debug("decoded block", "number", blockNumber, "hash", block.Hash().Hex())

	// Get CometBFT block by number (not hash, as Ethereum block hash may differ from CometBFT hash)
	resBlock, err := a.backend.CometBlockByNumber(blockNumber)
	if err != nil {
		a.logger.Debug("get block failed", "number", blockNumber, "error", err.Error())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		a.logger.Debug("block not found", "number", blockNumber)
		return nil, errors.New("block not found")
	}

	return a.backend.TraceBlock(blockNumber, config, resBlock)
}

// TraceCall lets you trace a given eth_call. It collects the structured logs
// created during the execution of EVM if the given transaction was added on
// top of the provided block and returns them as a JSON object.
func (a *API) TraceCall(args evmtypes.TransactionArgs, blockNrOrHash rpctypes.BlockNumberOrHash, config *rpctypes.TraceConfig) (interface{}, error) {
	a.logger.Debug("debug_traceCall", "args", args, "block number or hash", blockNrOrHash)
	return a.backend.TraceCall(args, blockNrOrHash, config)
}

// GetRawBlock retrieves the RLP-encoded block by block number or hash.
func (a *API) GetRawBlock(blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	a.logger.Debug("debug_getRawBlock", "block number or hash", blockNrOrHash)

	// Get block number from blockNrOrHash
	blockNum, err := a.backend.BlockNumberFromComet(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	// Get Ethereum block by number
	block, err := a.backend.EthBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("block not found")
	}

	// Encode block to RLP
	return rlp.EncodeToBytes(block)
}

// BlockProfile turns on goroutine profiling for nsec seconds and writes profile data to
// file. It uses a profile rate of 1 for most accurate information. If a different rate is
// desired, set the rate and write the profile manually.
func (a *API) BlockProfile(file string, nsec uint) error {
	a.logger.Debug("debug_blockProfile", "file", file, "nsec", nsec)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	runtime.SetBlockProfileRate(1)
	defer runtime.SetBlockProfileRate(0)

	time.Sleep(time.Duration(nsec) * time.Second) //#nosec G115 -- int overflow is not a concern here
	return writeProfile("block", file, a.logger)
}

// CpuProfile turns on CPU profiling for nsec seconds and writes
// profile data to file.
func (a *API) CpuProfile(file string, nsec uint) error { //nolint: revive
	a.logger.Debug("debug_cpuProfile", "file", file, "nsec", nsec)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	if err := a.StartCPUProfile(file); err != nil {
		return err
	}
	time.Sleep(time.Duration(nsec) * time.Second) //#nosec G115 -- int overflow is not a concern here
	return a.StopCPUProfile()
}

// GcStats returns GC statistics.
func (a *API) GcStats() (*debug.GCStats, error) {
	a.logger.Debug("debug_gcStats")
	if !a.profilingEnabled {
		return nil, rpctypes.ErrProfilingDisabled
	}
	s := new(debug.GCStats)
	debug.ReadGCStats(s)
	return s, nil
}

// GoTrace turns on tracing for nsec seconds and writes
// trace data to file.
func (a *API) GoTrace(file string, nsec uint) error {
	a.logger.Debug("debug_goTrace", "file", file, "nsec", nsec)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	if err := a.StartGoTrace(file); err != nil {
		return err
	}
	time.Sleep(time.Duration(nsec) * time.Second) //#nosec G115 -- int overflow is not a concern here
	return a.StopGoTrace()
}

// MemStats returns detailed runtime memory statistics.
func (a *API) MemStats() (*runtime.MemStats, error) {
	a.logger.Debug("debug_memStats")
	if !a.profilingEnabled {
		return nil, rpctypes.ErrProfilingDisabled
	}
	s := new(runtime.MemStats)
	runtime.ReadMemStats(s)
	return s, nil
}

// SetBlockProfileRate sets the rate of goroutine block profile data collection.
// rate 0 disables block profiling.
func (a *API) SetBlockProfileRate(rate int) error {
	a.logger.Debug("debug_setBlockProfileRate", "rate", rate)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	runtime.SetBlockProfileRate(rate)
	return nil
}

// Stacks returns a printed representation of the stacks of all goroutines.
func (a *API) Stacks() (string, error) {
	a.logger.Debug("debug_stacks")
	if !a.profilingEnabled {
		return "", rpctypes.ErrProfilingDisabled
	}
	buf := new(bytes.Buffer)
	err := pprof.Lookup("goroutine").WriteTo(buf, 2)
	if err != nil {
		a.logger.Error("Failed to create stacks", "error", err.Error())
	}
	return buf.String(), nil
}

// StartCPUProfile turns on CPU profiling, writing to the given file.
func (a *API) StartCPUProfile(file string) error {
	a.logger.Debug("debug_startCPUProfile", "file", file)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	a.handler.mu.Lock()
	defer a.handler.mu.Unlock()

	switch {
	case isCPUProfileConfigurationActivated(a.ctx):
		a.logger.Debug("CPU profiling already in progress using the configuration file")
		return errors.New("CPU profiling already in progress using the configuration file")
	case a.handler.cpuFile != nil:
		a.logger.Debug("CPU profiling already in progress")
		return errors.New("CPU profiling already in progress")
	default:
		fp, err := ExpandHome(file)
		if err != nil {
			a.logger.Debug("failed to get filepath for the CPU profile file", "error", err.Error())
			return err
		}
		f, err := os.Create(fp)
		if err != nil {
			a.logger.Debug("failed to create CPU profile file", "error", err.Error())
			return err
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			a.logger.Debug("cpu profiling already in use", "error", err.Error())
			if err := f.Close(); err != nil {
				a.logger.Debug("failed to close cpu profile file")
				return stderrors.Wrap(err, "failed to close cpu profile file")
			}
			return err
		}

		a.logger.Info("CPU profiling started", "profile", file)
		a.handler.cpuFile = f
		a.handler.cpuFilename = file
		return nil
	}
}

// StopCPUProfile stops an ongoing CPU profile.
func (a *API) StopCPUProfile() error {
	a.logger.Debug("debug_stopCPUProfile")
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	a.handler.mu.Lock()
	defer a.handler.mu.Unlock()

	switch {
	case isCPUProfileConfigurationActivated(a.ctx):
		a.logger.Debug("CPU profiling already in progress using the configuration file")
		return errors.New("CPU profiling already in progress using the configuration file")
	case a.handler.cpuFile != nil:
		a.logger.Info("Done writing CPU profile", "profile", a.handler.cpuFilename)
		pprof.StopCPUProfile()
		if err := a.handler.cpuFile.Close(); err != nil {
			a.logger.Debug("failed to close cpu file")
			return stderrors.Wrap(err, "failed to close cpu file")
		}
		a.handler.cpuFile = nil
		a.handler.cpuFilename = ""
		return nil
	default:
		a.logger.Debug("CPU profiling not in progress")
		return errors.New("CPU profiling not in progress")
	}
}

// WriteBlockProfile writes a goroutine blocking profile to the given file.
func (a *API) WriteBlockProfile(file string) error {
	a.logger.Debug("debug_writeBlockProfile", "file", file)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	return writeProfile("block", file, a.logger)
}

// WriteMemProfile writes an allocation profile to the given file.
// Note that the profiling rate cannot be set through the API,
// it must be set on the command line.
func (a *API) WriteMemProfile(file string) error {
	a.logger.Debug("debug_writeMemProfile", "file", file)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	return writeProfile("heap", file, a.logger)
}

// MutexProfile turns on mutex profiling for nsec seconds and writes profile data to file.
// It uses a profile rate of 1 for most accurate information. If a different rate is
// desired, set the rate and write the profile manually.
func (a *API) MutexProfile(file string, nsec uint) error {
	a.logger.Debug("debug_mutexProfile", "file", file, "nsec", nsec)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	runtime.SetMutexProfileFraction(1)
	time.Sleep(time.Duration(nsec) * time.Second) //#nosec G115 -- int overflow is not a concern here
	defer runtime.SetMutexProfileFraction(0)
	return writeProfile("mutex", file, a.logger)
}

// SetMutexProfileFraction sets the rate of mutex profiling.
func (a *API) SetMutexProfileFraction(rate int) error {
	a.logger.Debug("debug_setMutexProfileFraction", "rate", rate)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	runtime.SetMutexProfileFraction(rate)
	return nil
}

// WriteMutexProfile writes a goroutine blocking profile to the given file.
func (a *API) WriteMutexProfile(file string) error {
	a.logger.Debug("debug_writeMutexProfile", "file", file)
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	return writeProfile("mutex", file, a.logger)
}

// FreeOSMemory forces a garbage collection.
func (a *API) FreeOSMemory() error {
	a.logger.Debug("debug_freeOSMemory")
	if !a.profilingEnabled {
		return rpctypes.ErrProfilingDisabled
	}
	debug.FreeOSMemory()
	return nil
}

// SetGCPercent sets the garbage collection target percentage. It returns the previous
// setting. A negative value disables GC.
func (a *API) SetGCPercent(v int) (int, error) {
	a.logger.Debug("debug_setGCPercent", "percent", v)
	if !a.profilingEnabled {
		return 0, rpctypes.ErrProfilingDisabled
	}
	return debug.SetGCPercent(v), nil
}

// GetHeaderRlp retrieves the RLP encoded for of a single header.
func (a *API) GetHeaderRlp(number uint64) (hexutil.Bytes, error) {
	if !a.profilingEnabled {
		return nil, rpctypes.ErrProfilingDisabled
	}
	header, err := a.backend.HeaderByNumber(rpctypes.BlockNumber(number)) //#nosec G115 -- int overflow is not a concern here -- block number is not likely to exceed int64 max value
	if err != nil {
		return nil, err
	}

	return rlp.EncodeToBytes(header)
}

// GetBlockRlp retrieves the RLP encoded for of a single block.
func (a *API) GetBlockRlp(number uint64) (hexutil.Bytes, error) {
	if !a.profilingEnabled {
		return nil, rpctypes.ErrProfilingDisabled
	}
	block, err := a.backend.EthBlockByNumber(rpctypes.BlockNumber(number)) //#nosec G115 -- int overflow is not a concern here -- block number is not likely to exceed int64 max value
	if err != nil {
		return nil, err
	}

	return rlp.EncodeToBytes(block)
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (a *API) PrintBlock(number uint64) (string, error) {
	if !a.profilingEnabled {
		return "", rpctypes.ErrProfilingDisabled
	}
	block, err := a.backend.EthBlockByNumber(rpctypes.BlockNumber(number)) //#nosec G115 -- int overflow is not a concern here -- block number is not likely to exceed int64 max value
	if err != nil {
		return "", err
	}

	return spew.Sdump(block), nil
}

// IntermediateRoots executes a block, and returns a list
// of intermediate roots: the stateroot after each transaction.
func (a *API) IntermediateRoots(hash common.Hash, _ *evmtypes.TraceConfig) ([]common.Hash, error) {
	a.logger.Debug("debug_intermediateRoots", "hash", hash)
	if !a.profilingEnabled {
		return nil, rpctypes.ErrProfilingDisabled
	}
	return ([]common.Hash)(nil), nil
}
