package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"github.com/marcos-pereira-jr/simple-api-go/handlers"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			fx.Annotate(
        NewServeMux,
        fx.ParamTags(`group:"routes"`),
      ),
			AsRoute(handlers.NewEchoHandler),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

type Route interface {
	http.Handler

	Pattern() string
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
  for _, route := range routes {
    mux.Handle(route.Pattern(), route)
  }
	return mux
}

func AsRoute(f any) any {
  return fx.Annotate(
    f,
    fx.As(new(Route)),
    fx.ResultTags(`group:"routes"`),
    )
}

// https://uber-go.github.io/fx/get-started/echo-handler.html
func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
