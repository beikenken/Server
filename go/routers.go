/*
 * simple blog
 *
 * A Simple Blog
 *
 * API version: 1.0.0
 * Contact: apiteam@swagger.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/simpleblog/",
		Index,
	},

	Route{
		"DeleteArticleById",
		strings.ToUpper("Get"),
		"/simpleblog/user/deleteArticle/{id}",
		DeleteArticleById,
	},

	Route{
		"GetArticleById",
		strings.ToUpper("Get"),
		"/simpleblog/user/article/{id}",
		GetArticleById,
	},

	Route{
		"GetArticles",
		strings.ToUpper("Get"),
		"/simpleblog/user/articles",
		GetArticles,
	},

	Route{
		"CreateComment",
		strings.ToUpper("Post"),
		"/simpleblog/user/article/{id}/comments",
		CreateComment,
	},

	Route{
		"GetCommentsOfArticle",
		strings.ToUpper("Get"),
		"/simpleblog/user/article/{id}/comments",
		GetCommentsOfArticle,
	},

	Route{
		"SignIn",
		strings.ToUpper("Post"),
		"/simpleblog/user/signin",
		SignIn,
	},
}
