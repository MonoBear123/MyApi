package application

import (
	"net/http"

	"github.com/MonoBear123/MyApi/back/repository/expression"
	"github.com/MonoBear123/MyApi/server/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	//router.Handle("/front/*", http.StripPrefix("/front/", http.FileServer(http.Dir("/front/index.html"))))
	router.Route("/", a.loadExpressionRoutes)
	a.router = router

}
func (a *App) loadExpressionRoutes(router chi.Router) {
	ExpressionHandler := &handler.Expression{
		Repo: &expression.RedisRepo{
			Client: a.rdb,
		},
	}
	router.Post("/set", ExpressionHandler.SetExpression)
	router.Get("/get/{id}", ExpressionHandler.GetExpressionByID)
	router.Post("/setstatus", ExpressionHandler.SetAgentStatus)
	router.Get("/getstatus", ExpressionHandler.GetAgentStatus)

}
