package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/cqrs"
	"github.com/urfave/negroni"

	"github.com/danielkrainas/shex/api/server/handlers"
	"github.com/danielkrainas/shex/storage"
)

func New(ctx context.Context, config *configuration.Config) (*Server, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := acontext.GetLogger(ctx)
	log.Info("initializing server")

	storageDriver, err := storage.FromConfig(config)
	if err != nil {
		return nil, err
	}

	query := &cqrs.QueryDispatcher{
		Executors: []cqrs.QueryExecutor{
			storageDriver.Query(),
		},
	}

	command := &cqrs.CommandDispatcher{
		Handlers: []cqrs.CommandHandler{
			storageDriver.Command(),
		},
	}

	api, err := handlers.NewApi(query, command, config)
	if err != nil {
		return nil, fmt.Errorf("error creating server api: %v", err)
	}

	n := negroni.New()

	n.Use(cors.New(cors.Options{
		AllowedOrigins:   config.HTTP.CORS.Origins,
		AllowedMethods:   config.HTTP.CORS.Methods,
		AllowCredentials: true,
		AllowedHeaders:   config.HTTP.CORS.Headers,
		Debug:            config.HTTP.Debug,
	}))

	n.UseFunc(handlers.Logging)
	n.Use(handlers.Context(ctx))
	n.Use(&negroni.Recovery{
		Logger:     negroni.ALogger(log),
		PrintStack: true,
		StackAll:   true,
	})

	n.Use(handlers.Alive("/"))
	n.UseFunc(handlers.TrackErrors)
	n.UseHandler(api)

	s := &Server{
		Context: ctx,
		api:     api,
		config:  config,
		server: &http.Server{
			Addr:    config.HTTP.Addr,
			Handler: n,
		},
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	storage.LogSummary(ctx, config)

	return s, nil
}

type Server struct {
	context.Context
	config *configuration.Config
	server *http.Server
	api    *handlers.Api
}

func (server *Server) ListenAndServe() error {
	config := server.config
	ln, err := net.Listen("tcp", config.HTTP.Addr)
	if err != nil {
		return err
	}

	acontext.GetLogger(server.api).Infof("listening on %v", ln.Addr())
	return server.server.Serve(ln)
}

func configureLogging(ctx context.Context, config *configuration.Config) (context.Context, error) {
	log.SetLevel(logLevel(config.Log.Level))
	formatter := config.Log.Formatter
	if formatter == "" {
		formatter = "text"
	}

	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		if config.Log.Formatter != "" {
			return ctx, fmt.Errorf("unsupported log formatter: %q", config.Log.Formatter)
		}
	}

	if len(config.Log.Fields) > 0 {
		var fields []interface{}
		for k := range config.Log.Fields {
			fields = append(fields, k)
		}

		ctx = acontext.WithValues(ctx, config.Log.Fields)
		ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx, fields...))
	}

	ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level configuration.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
