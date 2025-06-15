// package tcp

// import (
// 	"context"
// 	"crypto/tls"
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"os"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/collectors"
// 	"github.com/prometheus/client_golang/prometheus/collectors/version"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/sirupsen/logrus"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// )

// type KokaqWireServerConfig struct {
// 	Port           int    `json:"port"`
// 	Timeout        int    `json:"timeout"`         // Timeout in seconds for each request
// 	UseTls         bool   `json:"use_tls"`         // Whether to use TLS for secure connections
// 	CertFile       string `json:"cert_file"`       // Path to the TLS certificate file
// 	KeyFile        string `json:"key_file"`        // Path to the TLS key file
// 	AuthMode       string `json:"auth_mode"`       // Authentication mode, e.g., "none", "basic", "token"
// 	MetricsAddress string `json:"metrics_address"` // Address for Prometheus metrics
// 	UsePrometheus  bool   `json:"use_prometheus"`  // Whether to enable Prometheus metrics
// 	MaxConnections int    `json:"max_connections"` // Maximum number of concurrent connections
// 	// Add more configuration options as needed
// }

// type KokaqWireServerMetrics struct {
// 	activeConnections prometheus.Gauge       `json:"active_connections"` // Gauge for active connections
// 	requestCounter    prometheus.Counter     `json:"request_counter"`    // Counter for total requests
// 	cpuTemp           prometheus.Gauge       `json:"cpu_temp"`           // Gauge for CPU temperature
// 	hdFailures        *prometheus.CounterVec `json:"hd_failures"`        // Counter for hard disk failures
// 	requestDurations  prometheus.Histogram   `json:"request_durations"`  // Histogram for request durations
// 	// Add more metrics as needed
// }

// type KokaqWireServer struct {
// 	port           int
// 	listener       net.Listener
// 	wg             sync.WaitGroup
// 	idleTimeout    int
// 	authMode       string
// 	useTls         bool
// 	certFile       string
// 	keyFile        string
// 	metricsAddress string
// 	metrics        *KokaqWireServerMetrics
// 	usePrometheus  bool
// 	logger         *logrus.Logger
// 	auditLogger    *logrus.Logger
// 	maxConnections int // Maximum number of concurrent connections
// }

// type IKokaqServerhandler interface {
// 	HandleRequest(*KokaqWireRequest) (*KokaqWireResponse, error)
// }

// // NewKokaqWireServer initializes a new KokaqWireServer with the provided port.
// func NewKokaqWireServer(config KokaqWireServerConfig) *KokaqWireServer {
// 	return &KokaqWireServer{
// 		port:           config.Port,
// 		idleTimeout:    config.Timeout,
// 		authMode:       config.AuthMode,
// 		useTls:         config.UseTls,
// 		certFile:       config.CertFile,
// 		keyFile:        config.KeyFile,
// 		logger:         logrus.New(),
// 		auditLogger:    logrus.New(),
// 		maxConnections: config.MaxConnections,
// 		metricsAddress: config.MetricsAddress,
// 		usePrometheus:  config.UsePrometheus,
// 		metrics: &KokaqWireServerMetrics{
// 			cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
// 				Name: "tcp_server_cpu_temperature",
// 				Help: "Current CPU temperature in Celsius",
// 			}),
// 			activeConnections: prometheus.NewGauge(prometheus.GaugeOpts{
// 				Name: "tcp_server_active_connections",
// 				Help: "Number of active TCP connections",
// 			}),
// 			requestCounter: prometheus.NewCounter(prometheus.CounterOpts{
// 				Name: "tcp_server_total_requests",
// 				Help: "Total number of processed requests",
// 			}),
// 			hdFailures: prometheus.NewCounterVec(
// 				prometheus.CounterOpts{
// 					Name: "tcp_server_hd_failures",
// 					Help: "Number of hard disk failures",
// 				},
// 				[]string{"disk"}),
// 			requestDurations: prometheus.NewHistogram(prometheus.HistogramOpts{
// 				Name:    "tcp_server_request_duration_seconds",
// 				Help:    "Duration of requests in seconds",
// 				Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5), // prometheus.DefBuckets, // Default buckets for request durations
// 			}),
// 		},
// 	}
// }

// func initTracer() func() {
// 	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		panic(err)
// 	}
// 	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
// 	otel.SetTracerProvider(tp)
// 	return func() {
// 		_ = tp.Shutdown(context.Background())
// 	}
// }

// // Start begins the TCP listener and handles incoming client requests.
// func (server *KokaqWireServer) Start(handler IKokaqServerhandler) error {

// 	var err error
// 	file, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
// 	if err != nil {
// 		server.logger.Fatalf("Could not create audit log: %v", err)
// 	}
// 	server.auditLogger.SetOutput(file)
// 	server.auditLogger.SetFormatter(&logrus.JSONFormatter{
// 		TimestampFormat: time.RFC3339,
// 	})
// 	server.auditLogger.SetLevel(logrus.InfoLevel)

// 	var reg = prometheus.NewRegistry()
// 	if server.usePrometheus {
// 		server.logger.Info("Starting Prometheus metrics server on ", server.metricsAddress)
// 		reg.MustRegister(
// 			collectors.NewGoCollector(),
// 			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
// 			version.NewCollector("example"),
// 			server.metrics.activeConnections,
// 			server.metrics.requestCounter,
// 			server.metrics.cpuTemp,
// 			server.metrics.hdFailures,
// 			server.metrics.requestDurations,
// 		)
// 		// Expose /metrics HTTP endpoint using the created custom registry.
// 		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
// 		// go http.ListenAndServe(server.metricsAddress, nil)
// 		cleanup := initTracer()
// 		defer cleanup()

// 	}

// 	var listener net.Listener
// 	address := fmt.Sprintf(":%d", server.port)
// 	if server.useTls {
// 		server.logger.Info("Starting server with TLS enabled.")
// 		var cer tls.Certificate
// 		cer, err = tls.LoadX509KeyPair(server.certFile, server.keyFile)
// 		if err != nil {
// 			server.logger.Error("failed to load TLS certs", err)
// 			return err
// 		}
// 		config := &tls.Config{Certificates: []tls.Certificate{cer}}
// 		listener, err = tls.Listen("tcp", address, config)
// 	} else {
// 		server.logger.Info("Starting server without TLS.")
// 		if server.useTls {
// 			server.logger.Warn("Warning: TLS is enabled but no certs provided, falling back to plain TCP.")
// 		}
// 		listener, err = net.Listen("tcp", address)
// 	}
// 	if err != nil {
// 		server.logger.Error("Error starting server ", err)
// 		return err
// 	}
// 	server.logger.WithField("addr", address).Info("TCP server started")
// 	server.listener = listener
// 	defer server.listener.Close()
// 	var connCount int32
// 	server.logger.Info("Listening on port ", server.port)
// 	for {
// 		server.logger.Info("Waiting for a connection...")
// 		client, err := server.listener.Accept()
// 		if err != nil {
// 			server.logger.WithError(err).Warn("Failed to accept connection")
// 			continue
// 		}
// 		server.logger.Info("Receiving a request...")

// 		// Check if the maximum number of connections is reached

// 		if atomic.AddInt32(&connCount, 1) > int32(server.maxConnections) {
// 			server.logger.Warn("Too many connections — rejecting client")
// 			client.Close()
// 			atomic.AddInt32(&connCount, -1)
// 			continue
// 		}

// 		// Handle each client in a separate goroutine
// 		server.wg.Add(1)
// 		go func(client net.Conn) {
// 			server.logger.Info("Received a request")
// 			defer server.wg.Done()
// 			defer client.Close()
// 			// Increment the active connections gauge
// 			server.metrics.activeConnections.Inc()
// 			defer server.metrics.activeConnections.Dec()

// 			tlsConn, ok := client.(*tls.Conn)
// 			if !ok {
// 				server.logger.Warn("Not a TLS connection")
// 				return
// 			}

// 			if err := tlsConn.Handshake(); err != nil {
// 				server.logger.WithError(err).Warn("TLS handshake failed")
// 				server.auditLogger.WithFields(logrus.Fields{
// 					"time":      time.Now().Format(time.RFC3339),
// 					"event":     "tls_handshake_failed",
// 					"error":     err.Error(),
// 					"remote_ip": client.RemoteAddr().String(),
// 				}).Warn("TLS handshake failed")
// 				return
// 			}

// 			// state := tlsConn.ConnectionState()

// 			// if len(state.PeerCertificates) == 0 {
// 			// 	server.logger.Warn("Client cert not provided")
// 			// 	return
// 			// }

// 			context := context.Background()

// 			tracer := otel.Tracer("kokaq-wire-server")
// 			context, span := tracer.Start(context, "HandleClientRequest")
// 			defer span.End()

// 			// Increment the request counter
// 			server.metrics.requestCounter.Inc()
// 			server.logger.WithField("remote_addr", client.RemoteAddr().String()).Info("Handling request from client")

// 			// Set a deadline for the client connection to avoid hanging indefinitely
// 			client.SetDeadline(time.Now().Add(time.Second * time.Duration(server.idleTimeout)))

// 			var req = &KokaqWireRequest{}

// 			err := req.ReadFromStream(client)
// 			if err != nil {
// 				server.logger.Error("Error reading request", err)
// 				return
// 			}

// 			switch req.MessageType {
// 			case MessageTypeOperational:
// 				server.logger.Info("Received an operational message.")
// 			case MessageTypeControl:
// 				server.logger.Info("Received a control message.")
// 			case MessageTypeAdmin:
// 				server.logger.Info("Received a adminops message.")
// 			default:
// 				server.logger.Info("Received an unknown message type.")
// 				return
// 			}

// 			switch req.OpCode {
// 			case OpCodeCreate:
// 				server.logger.Info("Received a create request.")
// 			case OpCodeDelete:
// 				server.logger.Info("Received a delete request.")
// 			case OpCodeGet:
// 				server.logger.Info("Received a get request.")
// 			case OpCodePeek:
// 				server.logger.Info("Received a peek request.")
// 			case OpCodePop:
// 				server.logger.Info("Received a pop request.")
// 			case OpCodePush:
// 				server.logger.Info("Received a push request.")
// 			case OpCodeAcquirePeekLock:
// 				server.logger.Info("Received an acquire peek lock request.")
// 			default:
// 				server.logger.Info("Received an unknown operation.")
// 				return
// 			}

// 			defer atomic.AddInt32(&connCount, -1)
// 			res, err := handler.HandleRequest(req)
// 			err = res.WriteToStream(client)
// 			if err != nil {
// 				server.logger.Error("Error writing response", err)
// 				return
// 			}
// 			server.logger.Info("Done handling the request.")
// 		}(client)
// 	}
// }

// // Stop stops the server and waits for all connections to be handled.
// func (server *KokaqWireServer) Stop() {
// 	server.listener.Close()
// 	server.wg.Wait()
// }
