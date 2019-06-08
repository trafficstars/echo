package fasthttpwithsocket

import (
	"github.com/labstack/echo/engine"
	fasthttpengine "github.com/trafficstars/echo/engine/fasthttp"
	"github.com/trafficstars/fasthttp"
	"github.com/trafficstars/fasthttpsocket"
)

type Server struct {
	*fasthttpengine.Server

	Socket *fasthttpsocket.Socket
}

func WithConfig(httpCfg engine.Config, sockCfg *fasthttpsocket.Config) (s *Server) {
	s = &Server{}
	s.Server = fasthttpengine.WithConfig(httpCfg)
	if sockCfg != nil {
		s.Socket = fasthttpsocket.NewSocket(s, *sockCfg)
	}
	return s
}

func (s *Server) SetHandler(h engine.Handler) {
	s.Server.SetHandler(h)
}

func (s *Server) Start() error {
	if s.Socket != nil {
		err := s.Socket.StartServer()
		if err != nil {
			return err
		}
	}

	err := s.Server.Start()
	if err != nil {
		s.Socket.StopServer()
	}
	return err
}

func (s *Server) Stop() error {
	err0 := s.Socket.StopServer()
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
