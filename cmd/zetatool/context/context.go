package context

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/clients"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
)

type Context struct {
	ctx            context.Context
	config         *config.Config
	zetacoreReader clients.ZetacoreReader
	inboundHash    string
	inboundChain   chains.Chain
	logger         zerolog.Logger
}

func NewContext(ctx context.Context, inboundChainID int64, inboundHash string, configFile string) (*Context, error) {
	observationChain, found := chains.GetChainFromChainID(inboundChainID, []chains.Chain{})
	if !found {
		return nil, fmt.Errorf("chain not supported,chain id: %d", inboundChainID)
	}
	cfg, err := config.GetConfig(observationChain, configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	zetacoreReader, err := clients.NewZetacoreReaderAdapter(cfg.ZetaChainRPC)
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        zerolog.Nop(),
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
	return &Context{
		ctx:            ctx,
		config:         cfg,
		zetacoreReader: zetacoreReader,
		inboundChain:   observationChain,
		inboundHash:    inboundHash,
		logger:         logger,
	}, nil
}

func (c *Context) GetContext() context.Context {
	return c.ctx
}

func (c *Context) GetConfig() *config.Config {
	return c.config
}

func (c *Context) GetZetacoreReader() clients.ZetacoreReader {
	return c.zetacoreReader
}

func (c *Context) GetInboundHash() string {
	return c.inboundHash
}

func (c *Context) GetInboundChain() chains.Chain {
	return c.inboundChain
}

func (c *Context) GetLogger() zerolog.Logger {
	return c.logger
}
