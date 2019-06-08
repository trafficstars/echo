package fasthttpwithsocket

import (
	"github.com/labstack/echo/engine"
	fasthttpengine "github.com/trafficstars/echo/engine/fasthttp"
	"github.com/trafficstars/fasthttp"
	"github.com/trafficstars/fasthttpsocket"
)

type Server struct {
	*fasthttpengine.Server

	Socket *fasthttpsocket.SocketServer
}

func New(httpAddress string, sockCfg *fasthttpsocket.Config) (s *Server, err error) {
	s = &Server{}
	s.Server = fasthttpengine.New(httpAddress)
	if sockCfg != nil {
		s.Socket, err = fasthttpsocket.NewSocketServer(s, *sockCfg)
	}
	return s, err
}

func WithConfig(httpCfg engine.Config, sockCfg *fasthttpsocket.Config) (s *Server, err error) {
	s = &Server{}
	s.Server = fasthttpengine.WithConfig(httpCfg)
	if sockCfg != nil {
		s.Socket, err = fasthttpsocket.NewSocketServer(s, *sockCfg)
	}
	return s, err
}

func (s *Server) SetHandler(h engine.Handler) {
	s.Server.SetHandler(h)
}

func (s *Server) Start() error {
	if s.Socket != nil {
		err := s.Socket.Start()
		if err != nil {
			return err
		}
	}

	err := s.Server.Start()
	if err != nil {
		s.Socket.Stop()
	}
	return err
}

func (s *Server) Stop() error {
	err0 := s.Socket.Stop()
	err1 := s.Server.Stop()
	if err1 != nil {
		return err1
	}
	return err0
}

func (s *Server) HandleRequest(ctx *fasthttp.RequestCtx) error {
	s.Handler(ctx)
	return nil
}
