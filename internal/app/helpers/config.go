package helpers

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Options struct {
	ServerAddress      string `env:"SERVER_ADDRESS"`
	DBConnectionString string `env:"DATABASE_DSN"`
	SecretKey          string `env:"SECRET_KEY"`
}

var instance *Options

func ReadFlags() *Options {
	if instance == nil {
		cmdOptions := getCmdOptions()
		envOptions := getEnvOptions()

		finalOptions := Options{}
		// env options are the priority
		mergeOptions(&finalOptions, envOptions)
		mergeOptions(&finalOptions, cmdOptions)
		instance = &finalOptions
	}
	return instance
}

func mergeOptions(mergeInto *Options, newValues Options) {
	if mergeInto.ServerAddress == "" && newValues.ServerAddress != "" {
		mergeInto.ServerAddress = newValues.ServerAddress
	}

	if mergeInto.ShortenerBaseURL == "" && newValues.ShortenerBaseURL != "" {
		mergeInto.ShortenerBaseURL = newValues.ShortenerBaseURL
	}

	if mergeInto.FileStoragePath == "" && newValues.FileStoragePath != "" {
		mergeInto.FileStoragePath = newValues.FileStoragePath
	}

	if mergeInto.DBConnectionString == "" && newValues.DBConnectionString != "" {
		mergeInto.DBConnectionString = newValues.DBConnectionString
	}

	if mergeInto.SecretKey == "" && newValues.SecretKey != "" {
		mergeInto.SecretKey = newValues.SecretKey
	}
}

func getEnvOptions() Options {
	var opt Options
	err := env.Parse(&opt)
	if err != nil {
		log.Fatalln(err)
	}
	return opt
}

func getCmdOptions() Options {
	opt := Options{}
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080", "port on which the server should run")
	flag.StringVar(&opt.DBConnectionString, "d",
		"",
		"Postgres database connection string")
	flag.StringVar(&opt.SecretKey, "s",
		"supersecretkey",
		"Secret key for data encryption")
	flag.Parse()
	return opt
}
