package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/abassian/huron/src/node"
	"github.com/abassian/huron/src/peers"
	"github.com/sirupsen/logrus"
)

// Service ...
type Service struct {
	bindAddress string
	node        *node.Node
	graph       *node.Graph
	logger      *logrus.Logger
}

// NewService ...
func NewService(bindAddress string, n *node.Node, logger *logrus.Logger) *Service {
	service := Service{
		bindAddress: bindAddress,
		node:        n,
		graph:       node.NewGraph(n),
		logger:      logger,
	}

	return &service
}

//Serve defines endpoints and starts ListenAndServe
func (s *Service) Serve() {
	s.logger.WithField("bind_address", s.bindAddress).Debug("Huron Service serving")

	serverMuxHuron := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/stats", s.GetStats)
	r.HandleFunc("/block/{index}", s.GetBlock)
	r.HandleFunc("/graph", s.GetGraph)
	r.HandleFunc("/peers", s.GetPeers)
	r.HandleFunc("/genesispeers", s.GetGenesisPeers)

	serverMuxHuron.Handle("/", &CORSServer{r})

	err := http.ListenAndServe(s.bindAddress, serverMuxHuron)

	if err != nil {
		s.logger.WithField("error", err).Error("Service failed")
	}
}

// CORSServer ...
type CORSServer struct {
	r *mux.Router
}

// ServeHTTP ...
func (s *CORSServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

// GetStats ...
func (s *Service) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := s.node.GetStats()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stats)
}

// GetBlock ...
func (s *Service) GetBlock(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/block/"):]

	blockIndex, err := strconv.Atoi(param)

	if err != nil {
		s.logger.WithError(err).Errorf("Parsing block_index parameter %s", param)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	block, err := s.node.GetBlock(blockIndex)

	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving block %d", blockIndex)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(block)
}

// GetGraph ...
func (s *Service) GetGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	res, _ := s.graph.GetInfos()

	encoder.Encode(res)
}

// GetPeers ...
func (s *Service) GetPeers(w http.ResponseWriter, r *http.Request) {
	returnPeerSet(w, r, s.node.GetPeers())
}

// GetGenesisPeers ...
func (s *Service) GetGenesisPeers(w http.ResponseWriter, r *http.Request) {
	returnPeerSet(w, r, s.node.GetGenesisPeers())
}

func returnPeerSet(w http.ResponseWriter, r *http.Request, peers []*peers.Peer) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	encoder.Encode(peers)
}
