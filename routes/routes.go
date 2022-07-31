package routes

import (
	"hexagony/app/domain"
	controller "hexagony/app/http/controllers"
	"hexagony/app/http/middleware"

	"github.com/go-chi/chi/v5"
)

type RoutesUseCases struct {
	domain.AuthUseCase
	domain.UsersUseCase
	domain.AlbumsUseCase
}

func AuthRoutes(c *chi.Mux, auc domain.AuthUseCase) {
	handler := controller.AuthController{AuthUseCase: auc}

	c.Post("/auth", handler.Authenticate)
}

func UsersRoutes(c *chi.Mux, as domain.UsersUseCase) {
	handler := controller.UsersController{UsersUseCase: as}

	c.Route("/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/", handler.FindAll)
		r.Get("/{uuid}", handler.FindByID)
		r.Post("/", handler.Add)
		r.Put("/{uuid}", handler.Update)
		r.Delete("/{uuid}", handler.Delete)
	})
}

func AlbumsRoutes(c *chi.Mux, as domain.AlbumsUseCase) {
	handler := controller.AlbumsController{AlbumsUseCase: as}

	c.Route("/album", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/", handler.FindAll)
		r.Get("/{uuid}", handler.FindByID)
		r.Post("/", handler.Add)
		r.Put("/{uuid}", handler.Update)
		r.Delete("/{uuid}", handler.Delete)
	})
}

func Api(c *chi.Mux, r *RoutesUseCases) {
	AuthRoutes(c, r.AuthUseCase)
	UsersRoutes(c, r.UsersUseCase)
	AlbumsRoutes(c, r.AlbumsUseCase)
}
