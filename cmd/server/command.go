package main

import (
	"flag"
	"godrop/internal/client"
	"godrop/internal/common"
	"godrop/internal/common/log"
	"godrop/internal/server"
	"log/slog"
	"os"
)

type Command interface {
	Name() string
	Describe() string
	Execute()
}

type ServerCommand struct {
}

func (s *ServerCommand) Name() string {
	return "serve"
}

func (s *ServerCommand) Describe() string {
	return "used for serving file or directory"
}

func (s *ServerCommand) Execute() {
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	var (
		PathToServe = ""
		Password    = ""
		QrEnabled   = false
		NTimes      = uint(1)
		Port        = uint(8080)
		LogLevel    = ""
	)
	serveCmd.StringVar(&PathToServe, "path", "", "Path to specific folder or file to share")
	serveCmd.StringVar(&Password, "password", common.DefaultPassword, "Password to protected from intruders")
	serveCmd.BoolVar(&QrEnabled, "qr", server.DefaultQrOptions, "Enable QR code sharing")
	serveCmd.UintVar(&NTimes, "times", server.DefaultTimes, "Define how many times file can be downloaded before closing")
	serveCmd.UintVar(&Port, "port", uint(server.DefaultPort), "Port for server")
	serveCmd.StringVar(&LogLevel, "log-level", log.DefaultLogLevelStr, "Logging level one of [DEBUG,INFO,WARN,ERROR]")

	err := serveCmd.Parse(os.Args[2:])
	if err != nil {
		slog.Error("Error while parsing arguments", "error", err)
		os.Exit(1)
	}
	port := uint16(Port)
	err = server.Start(server.ServerOptions{
		PathToServe: &PathToServe,
		Password:    &Password,
		QrEnabled:   QrEnabled,
		NTimes:      &NTimes,
		Port:        &port,
		LogLevel:    &LogLevel,
	})
	if err != nil {
		slog.Error("Error while starting server", "error", err)
		os.Exit(1)
	}
}

type GetCommand struct {
}

func (g *GetCommand) Name() string {
	return "get"
}

func (g *GetCommand) Describe() string {
	return "used for downloading file or directory from host"
}

func (g *GetCommand) Execute() {
	var (
		From         = ""
		UserPassword = ""
		RequestPort  = uint(server.DefaultPort)
	)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getCmd.StringVar(&From, "from", "", "IP address to get file from")
	getCmd.StringVar(&UserPassword, "password", common.DefaultPassword, "Password for getting resource")
	getCmd.UintVar(&RequestPort, "port", uint(server.DefaultPort), "Default port for remote address")

	err := getCmd.Parse(os.Args[2:])
	if err != nil {
		slog.Error("Error while parsing arguments", "error", err)
		os.Exit(1)
	}

	err = client.GetFile(client.ClientOptions{
		From:     &From,
		Password: &UserPassword,
		Port:     RequestPort,
	})

	if err != nil {
		slog.Error("Error while getting file", "error", err)
		os.Exit(1)
	}

}
