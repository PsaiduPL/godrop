package server

import (
	"context"
	"fmt"
	"godrop/internal/common"
	"godrop/internal/server/file"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
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
}

// postvalidated Server Config
type ServerConfig struct {
	PathToServer string
	Password     string
	QrEnabled    bool
	NTimes       uint
	Port         uint16
	ShareFile    file.ShareFile
}

type DropServer struct {
	ServerConfig ServerConfig

	Counter     uint
	CounterLock sync.RWMutex
}

func (opts *ServerOptions) validate() (*ServerConfig, error) {
	var (
		password  = common.DefaultPassword
		port      = DefaultPort
		times     = DefaultTimes
		qrEnabled = DefaultQrOptions
		file      *file.ShareFile
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

	return &ServerConfig{
		Password:  password,
		Port:      port,
		QrEnabled: qrEnabled,
		NTimes:    times,
		ShareFile: *file,
	}, nil
}
func validatePath(path *string) (*file.ShareFile, error) {
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

func (dropServer *DropServer) start() error {
	informAboutIp()
	config := dropServer.ServerConfig
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
	mux.HandleFunc("GET /api/v1/share", HandleFileRequest)
	error := server.ListenAndServe()
	if error != nil {
		return error
	}
	defer server.Close()
	return nil
}

func HandleFileRequest(responseWriter http.ResponseWriter, req *http.Request) {
	slog.Info("New request received")
	server := req.Context().Value("server").(*DropServer)
	config := server.ServerConfig
	hashedPasswordHeader := req.Header.Get(common.HeaderName)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordHeader), []byte(config.Password)); err != nil {
		slog.Warn("Invalid password for resource")
		http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
		return
	}
	responseWriter.Header().Add("Content-Type", "application/octet-stream")
	fileName := path.Base(config.ShareFile.Path)

	responseWriter.Header().Add("X-FileName", fileName)
	ext := ""
	ext = path.Ext(config.ShareFile.Path)
	if config.ShareFile.IsDir {
		ext = ".zip"
	}
	responseWriter.Header().Add("X-FileExtension", ext)
	responseWriter.WriteHeader(http.StatusOK)
	err := config.ShareFile.WriteContent(responseWriter)
	if err != nil {
		slog.Error("Error while writing content", "error", err)
		return
	}
	server.CounterLock.Lock()
	defer server.CounterLock.Unlock()
	server.Counter -= 1
	slog.Info("Succesfuly downloaded file ", "counter", server.Counter)
	if server.Counter == 0 {
		slog.Info("Exiting program as it was succesfuly downloaded by all users")
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)

		}()
	}

}

func informAboutIp() {
	ip, err := getLocalIPs()
	if err != nil {
		slog.Error("Error while getting local address", "error", err)
		os.Exit(1)
	}

	slog.Info("Current local ip address", "ip", ip[0])
}

func getLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("No addreses found")
	}
	return ips, nil
}
