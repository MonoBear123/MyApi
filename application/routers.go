package application

import (
	"net/http"

	"github.com/MonoBear123/MyApi/handler"
	"github.com/MonoBear123/MyApi/repository/expression"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/ex", a.loadExpressionRoutes)
	a.router = router
}
func (a *App) loadExpressionRoutes(router chi.Router) {
	ExpressionHandler := &handler.Expression{
		Repo: &expression.RedisRepo{
			Client: a.rdb,
		},
	}
	router.Post("/", ExpressionHandler.SetExpression)
	router.Get("/{id}", ExpressionHandler.GetExpressionByID)

}
