		package main

		import (
			"fmt"
			"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
		)

		func main() {
			c := config.ServiceConfig{
				Name: "test",
				BaseURL: "http://localhost:3000",
				Routes: []config.RouteConfig{
					{
						Path: "/ping",
						Methods: []string{"GET"},
						AuthRequired: false,
					},
				},
			}

			fmt.Printf("%+v\n", c)
		}
