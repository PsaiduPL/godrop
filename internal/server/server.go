package server

import (
	"context"
	"fmt"
	"godrop/internal/common"
	"godrop/internal/common/colored"
	"godrop/internal/common/log"
	"godrop/internal/server/file"
	"godrop/internal/server/ip"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultPort      = uint16(9999)
	DefaultTimes     = uint(1)
	DefaultQrOptions = false
)

// prevalidated options
type ServerOptions struct {
	PathToServe *string
	Password    *string
	QrEnabled   bool
	NTimes      *uint
	Port        *uint16
	LogLevel    *string
}

// postvalidated Server Config
type ServerConfig struct {
	PathToServer string
	Password     string
	QrEnabled    bool
	NTimes       uint
	Port         uint16
	ShareFile    file.ShareFile
	LogLevel     slog.Level
}

type DropServer struct {
	ServerConfig ServerConfig

	Counter     uint
	CounterLock sync.RWMutex
}

func Start(options ServerOptions) error {
	config, err := options.validate()
	if err != nil {
		return err
	}
	server := DropServer{
		ServerConfig: *config,
		Counter:      config.NTimes,
	}
	return server.start()
}

func (opts *ServerOptions) validate() (*ServerConfig, error) {
	var (
		password  = common.DefaultPassword
		port      = DefaultPort
		times     = DefaultTimes
		qrEnabled = DefaultQrOptions
		file      *file.ShareFile
		logLevel  = log.DefaultLogLevel
	)
	if opts.Password != nil {
		password = *opts.Password
	}
	if opts.NTimes != nil {
		times = *opts.NTimes
	}
	if !opts.QrEnabled {
		qrEnabled = opts.QrEnabled
	}
	if opts.Port != nil {
		port = *opts.Port
	}
	file, err := validatePath(opts.PathToServe)
	if err != nil {
		return nil, err
	}
	if opts.LogLevel != nil {
		level, err := log.MapStringToSlogLevel(*opts.LogLevel)
		if err != nil {
			return nil, err
		}
		logLevel = level

	}

	return &ServerConfig{
		Password:  password,
		Port:      port,
		QrEnabled: qrEnabled,
		NTimes:    times,
		ShareFile: *file,
		LogLevel:  logLevel,
	}, nil
}
func validatePath(path *string) (*file.ShareFile, error) {
	if path == nil || strings.TrimSpace(*path) == "" {
		return nil, fmt.Errorf("No path provided")
	}
	var (
		stat  os.FileInfo
		isDir = false
		err   error
	)

	if stat, err = os.Stat(*path); err != nil && os.IsExist(err) {
		return nil, err
	}
	if stat.IsDir() {
		isDir = true
	}

	return &file.ShareFile{
		Path:  *path,
		IsDir: isDir,
	}, nil

}

func (dropServer *DropServer) start() error {
	ip.InformAboutIp()
	config := dropServer.ServerConfig
	log.SetLogLevel(config.LogLevel)
	mux := http.NewServeMux()
	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", config.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "server", dropServer)
			return ctx
		},
	}

	slog.Info("Server started ", "port", config.Port)
	mux.HandleFunc("GET /api/v1/share", handleFileRequest)
	error := server.ListenAndServe()
	if error != nil {
		return error
	}
	defer server.Close()
	return nil
}

func handleFileRequest(responseWriter http.ResponseWriter, req *http.Request) {
	colored.PrintColoredWithTags("<GREEN>New request received<GREEN>\n")
	// slog.Info("New request received")
	var (
		server = req.Context().Value("server").(*DropServer)
		config = server.ServerConfig
		err    error
	)
	err = checkHashInHeader(&config, &req.Header)
	if err != nil {
		http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
	}

	writeContentTypeAndFileName(&config, responseWriter)
	err = config.ShareFile.WriteContent(responseWriter)
	if err != nil {
		slog.Error("Error while writing content", "error", err)
		return
	}
	server.CounterLock.Lock()
	defer server.CounterLock.Unlock()
	server.Counter -= 1
	slog.Info("Succesfuly downloaded file ", "counter", server.Counter)
	if server.Counter == 0 {
		colored.PrintColoredWithTags("<RED>Exiting program as it was succesfuly downloaded by all users\n<RED>")
		// slog.Info("Exiting program as it was succesfuly downloaded by all users")
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)

		}()
	}

}

func checkHashInHeader(serverConfig *ServerConfig, header *http.Header) error {
	hashedPasswordHeader := header.Get(common.HeaderName)
	if strings.TrimSpace(hashedPasswordHeader) == "" {
		return fmt.Errorf("Empty header with password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordHeader), []byte(serverConfig.Password)); err != nil {
		slog.Warn("Invalid password for resource")
		return fmt.Errorf("Invalid password")
	}
	return nil
}

func writeContentTypeAndFileName(serverConfig *ServerConfig, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/octet-stream")
	fileName := filepath.Base(filepath.Clean(serverConfig.ShareFile.Path)) // if its folder add zip extension
	if serverConfig.ShareFile.IsDir {
		fileName += ".zip"
	}
	writer.Header().Add("X-FileName", fileName)
	ext := ""
	ext = filepath.Ext(serverConfig.ShareFile.Path)
	if serverConfig.ShareFile.IsDir {
		ext = ".zip"
	}
	writer.Header().Add("X-FileExtension", ext)
	writer.WriteHeader(http.StatusOK)

}
