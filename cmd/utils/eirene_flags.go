package utils

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zenanet-network/go-zenanet/eth"
	"github.com/zenanet-network/go-zenanet/eth/ethconfig"
	"github.com/zenanet-network/go-zenanet/node"
)

var (
	//
	// Eirene Specific flags
	//

	// HarmoniaURLFlag flag for harmonia url
	HarmoniaURLFlag = &cli.StringFlag{
		Name:  "eirene.harmonia",
		Usage: "URL of Harmonia service",
		Value: "http://localhost:1317",
	}

	// WithoutHarmoniaFlag no harmonia (for testing purpose)
	WithoutHarmoniaFlag = &cli.BoolFlag{
		Name:  "eirene.withoutharmonia",
		Usage: "Run without Harmonia service (for testing purpose)",
	}

	// HarmoniagRPCAddressFlag flag for harmonia gRPC address
	HarmoniagRPCAddressFlag = &cli.StringFlag{
		Name:  "eirene.harmoniagRPC",
		Usage: "Address of Harmonia gRPC service",
		Value: "",
	}

	// RunHarmoniaFlag flag for running harmonia internally from eirene
	RunHarmoniaFlag = &cli.BoolFlag{
		Name:  "eirene.runharmonia",
		Usage: "Run Harmonia service as a child process",
	}

	RunHarmoniaArgsFlag = &cli.StringFlag{
		Name:  "eirene.runharmoniaargs",
		Usage: "Arguments to pass to Harmonia service",
		Value: "",
	}

	// UseHarmoniaApp flag for using internal harmonia app to fetch data
	UseHarmoniaAppFlag = &cli.BoolFlag{
		Name:  "eirene.useharmoniaapp",
		Usage: "Use child harmonia process to fetch data, Only works when eirene.runharmonia is true",
	}

	// EireneFlags all eirene related flags
	EireneFlags = []cli.Flag{
		HarmoniaURLFlag,
		WithoutHarmoniaFlag,
		HarmoniagRPCAddressFlag,
		RunHarmoniaFlag,
		RunHarmoniaArgsFlag,
		UseHarmoniaAppFlag,
	}
)

// SetEireneConfig sets eirene config
func SetEireneConfig(ctx *cli.Context, cfg *eth.Config) {
	cfg.HarmoniaURL = ctx.String(HarmoniaURLFlag.Name)
	cfg.WithoutHarmonia = ctx.Bool(WithoutHarmoniaFlag.Name)
	cfg.HarmoniagRPCAddress = ctx.String(HarmoniagRPCAddressFlag.Name)
	cfg.RunHarmonia = ctx.Bool(RunHarmoniaFlag.Name)
	cfg.RunHarmoniaArgs = ctx.String(RunHarmoniaArgsFlag.Name)
	cfg.UseHarmoniaApp = ctx.Bool(UseHarmoniaAppFlag.Name)
}

// CreateEireneZenanet Creates eirene zenanet object from eth.Config
func CreateEireneZenanet(cfg *ethconfig.Config) *eth.Zenanet {
	workspace, err := os.MkdirTemp("", "eirene-command-node-")
	if err != nil {
		Fatalf("Failed to create temporary keystore: %v", err)
	}

	// Create a networkless protocol stack and start an Zenanet service within
	stack, err := node.New(&node.Config{DataDir: workspace, UseLightweightKDF: true, Name: "eirene-command-node"})
	if err != nil {
		Fatalf("Failed to create node: %v", err)
	}

	zenanet, err := eth.New(stack, cfg)
	if err != nil {
		Fatalf("Failed to register Zenanet protocol: %v", err)
	}

	// Start the node and assemble the JavaScript console around it
	if err = stack.Start(); err != nil {
		Fatalf("Failed to start stack: %v", err)
	}

	stack.Attach()

	return zenanet
}
