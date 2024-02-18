package application

import (
	"html/template"
	"log"
	"net/http"

	"github.com/MonoBear123/MyApi/back/repository/expression"
	"github.com/MonoBear123/MyApi/server/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/", a.loadExpressionRoutes)
	a.router = router

}
func (a *App) loadExpressionRoutes(router chi.Router) {
	ExpressionHandler := &handler.Expression{
		Repo: &expression.RedisRepo{
			Client: a.rdb,
		},
	}
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("front/index.html")
		var expr string
		err := t.Execute(w, expr)
		if err != nil {
			log.Println("error")
		}
	})

	router.Post("/set", ExpressionHandler.SetExpression)

	router.Post("/setstatus", ExpressionHandler.SetAgentStatus)
	router.Get("/getstatus", ExpressionHandler.GetAgentStatus)
	router.Get("/config", ExpressionHandler.UpdateConfigHandler)

}
