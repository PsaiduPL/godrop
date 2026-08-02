package main

import (
	"flag"
	"godrop/internal/client"
	"godrop/internal/common"
	"godrop/internal/server"
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("No args provided")
		os.Exit(1)
	}
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	var (
		PathToServe = ""
		Password    = ""
		QrEnabled   = false
		NTimes      = uint(1)
		Port        = uint(8080)
	)
	serveCmd.StringVar(&PathToServe, "path", "", "Path to specific folder or file to share")
	serveCmd.StringVar(&Password, "password", common.DefaultPassword, "Password to protected from intruders")
	serveCmd.BoolVar(&QrEnabled, "qr", server.DefaultQrOptions, "Enable QR code sharing")
	serveCmd.UintVar(&NTimes, "times", server.DefaultTimes, "Define how many times file can be downloaded before closing")
	serveCmd.UintVar(&Port, "port", uint(server.DefaultPort), "Port for server")

	var (
		From         = ""
		UserPassword = ""
		RequestPort  = uint(server.DefaultPort)
	)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getCmd.StringVar(&From, "from", "", "IP address to get file from")
	getCmd.StringVar(&UserPassword, "password", common.DefaultPassword, "Password for getting resource")
	getCmd.UintVar(&RequestPort, "port", uint(server.DefaultPort), "Default port for remote address")
	switch os.Args[1] {
	case "serve":
		{
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
			})
			if err != nil {
				slog.Error("Error while starting server", "error", err)
				os.Exit(1)
			}

		}
	case "get":
		{
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
				slog.Error("Error while starting server", "error", err)
				os.Exit(1)
			}

		}
	}

}
