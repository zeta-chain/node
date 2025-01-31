package context

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/rpc"
	zetacorerpc "github.com/zeta-chain/node/pkg/rpc"
)

type Context struct {
	ctx            context.Context
	config         *config.Config
	zetaCoreClient rpc.Clients
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
	zetacoreClient, err := zetacorerpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to create zetacore client: %w", err)
	}
	// logger is used when calling internal zetaclient functions which need a logger.
	// we do not need to log those messages for this tool
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        zerolog.Nop(),
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
	return &Context{
		ctx:            ctx,
		config:         cfg,
		zetaCoreClient: zetacoreClient,
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

func (c *Context) GetZetaCoreClient() rpc.Clients {
	return c.zetaCoreClient
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
