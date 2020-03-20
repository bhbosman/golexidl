package main

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"io"
	"log"
	"os"
)

var (
	stdin  = bufio.NewReader(os.Stdin)
	stdout = bufio.NewWriter(os.Stdout)
	stderr = bufio.NewWriter(os.Stderr)
)

func main() {
	flags := NewAppCtx()
	flags.run()

	getLogger := func(verbose bool) io.Writer {
		if verbose {
			return os.Stdout
		}
		return &nullWriter{}
	}

	logger := log.New(getLogger(flags.verbose), "idllexgo: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	var factory processFactory
	//var process *processLexemsImpl
	app := fx.New(
		fx.Logger(&loggerWrapper{
			applicationLogger: logger,
		}),
		fx.Provide(func() *log.Logger {
			return logger
		}),
		fx.Provide(func(appCtx *GoLexIdlAppCtx, logger *log.Logger) processFactory {
			return &processFactoryImpl{
				ctx:    appCtx,
				logger: logger,
			}
		}),
		fx.Provide(
			func(lc fx.Lifecycle) *GoLexIdlAppCtx {
				return flags
			}),
		fx.Provide(
			func(flags *GoLexIdlAppCtx) *ResolveInputFileName {
				return &ResolveInputFileName{}
			}),
		fx.Invoke(func(lc fx.Lifecycle, item *ResolveInputFileName, flags *GoLexIdlAppCtx) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					var err error
					flags.inputFileName, err = item.run(flags.inputFileName)
					return err
				},
				OnStop: nil,
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, flags *GoLexIdlAppCtx) {
			var l *os.File
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					var err error
					if flags.inputFileName != "" {
						stat, err := os.Stat(flags.inputFileName)
						if err != nil {
							return err
						}
						if stat.IsDir() {
							return fmt.Errorf("can not open folder")
						}
					}
					if flags.inputFileName == "" {
						flags.inputReader = stdin
						flags.inputFileName = "(stdin)"
					} else {
						l, err = os.Open(flags.inputFileName)
						if err != nil {
							return err
						}
						flags.inputReader = bufio.NewReader(l)
					}
					return err
				},
				OnStop: func(ctx context.Context) error {
					if l != nil {
						return l.Close()
					}
					return nil
				},
			})
		}),

		fx.Invoke(func(lc fx.Lifecycle, flags *GoLexIdlAppCtx) {
			var g *os.File
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					var err error
					if flags.tflag {
						flags.outputWriter = stdout
					} else {
						if flags.oflag == "" {
							log.Fatal(fmt.Errorf("no output file"))
						}
						g, err = os.Create(flags.oflag)
						if err != nil {
							log.Fatal(err)
						}
						flags.outputWriter = bufio.NewWriterSize(g, 1024)
					}
					return err
				},
				OnStop: func(ctx context.Context) error {
					err := flags.outputWriter.Flush()
					if g != nil {
						err = multierr.Append(err, g.Close())
					}
					return err
				},
			})
		}),
		fx.Populate(&factory))

	startError := app.Start(context.TODO())
	if startError != nil {
		os.Exit(1)
	}
	defer func() {
		_ = app.Stop(context.TODO())
	}()
	process := factory.create()
	process.run()
}
