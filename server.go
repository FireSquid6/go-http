package ink_server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type RequestHandler func(req Request) Response

type Request struct {
	Method string
	Path   string
	Body   string
}

type Response struct {
	Status int
	Body   string
}

type Route struct {
	Method  string
	Path    string
	Handler RequestHandler
}

type Server struct {
	routes []Route
}

// handler functions
func (s *Server) Get(path string, handler RequestHandler) {
	s.routes = append(s.routes, Route{Method: "GET", Path: path, Handler: handler})
}

func (s *Server) Post(path string, handler RequestHandler) {
	s.routes = append(s.routes, Route{Method: "POST", Path: path, Handler: handler})
}

func (s *Server) Put(path string, handler RequestHandler) {
	s.routes = append(s.routes, Route{Method: "PUT", Path: path, Handler: handler})
}

func (s *Server) Delete(path string, handler RequestHandler) {
	s.routes = append(s.routes, Route{Method: "DELETE", Path: path, Handler: handler})
}

func (s *Server) Patch(path string, handler RequestHandler) {
	s.routes = append(s.routes, Route{Method: "PATCH", Path: path, Handler: handler})
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
      req := getRequest(c)
      res := s.handle(req)

      c.Write([]byte(responseToString(res)))
			c.Close()
		}(c)
	}
}

func responseToString(res Response) string {
  return fmt.Sprintf("HTTP/1.1 %d OK\r\n\r\n%s", res.Status, res.Body)
}

func (s *Server) handle(req Request) Response {
  for _, route := range s.routes {
    if route.Method == req.Method && route.Path == req.Path {
      return route.Handler(req)
    }
  }

  return Response{Status: 404, Body: "Not Found"}
}

func getRequest(c net.Conn) Request {
  buffer := make([]byte, 1024)
  c.Read(buffer)

  request := string(buffer)
  fmt.Println("Received request:", request)

  return Request{}
}

func handleConnection(c net.Conn) {
	buffer := make([]byte, 1024)
	c.Read(buffer)

	request := string(buffer)
	fmt.Println("Received request:", request)

	if strings.HasPrefix(request, "GET / HTTP/1.1") {
		c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
