package servers

import (
	"github.com/gorilla/mux"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications"
)

func New(
	application applications.Application,
) *mux.Router {
	api := &api{
		application: application,
	}

	router := mux.NewRouter()

	router.HandleFunc("/platforms", api.createPlatform).Methods("POST")
	router.HandleFunc("/platforms", api.findPlatforms).Methods("GET")
	router.HandleFunc("/platforms/count", api.countPlatforms).Methods("GET")
	router.HandleFunc("/platforms/{id}", api.findPlatformByID).Methods("GET")
	router.HandleFunc("/platforms/handle/{handle}", api.findPlatformByHandle).Methods("GET")

	router.HandleFunc("/users", api.createUser).Methods("POST")
	router.HandleFunc("/users", api.findUsers).Methods("GET")
	router.HandleFunc("/users/count", api.countUsers).Methods("GET")
	router.HandleFunc("/users/{id}", api.findUserByID).Methods("GET")
	router.HandleFunc("/users/platform/{platform_id}/external/{external_id}", api.findUserByExternalID).Methods("GET")
	router.HandleFunc("/users/platform/{platform_id}/handle/{handle}", api.findUserByHandle).Methods("GET")

	router.HandleFunc("/communities", api.createCommunity).Methods("POST")
	router.HandleFunc("/communities", api.findCommunities).Methods("GET")
	router.HandleFunc("/communities/count", api.countCommunities).Methods("GET")
	router.HandleFunc("/communities/{id}", api.findCommunityByID).Methods("GET")
	router.HandleFunc("/communities/platform/{platform_id}", api.findCommunitiesByPlatform).Methods("GET")
	router.HandleFunc("/communities/platform/{platform_id}/handle/{handle}", api.findCommunityByHandle).Methods("GET")

	router.HandleFunc("/posts", api.createPost).Methods("POST")
	router.HandleFunc("/posts", api.findPosts).Methods("GET")
	router.HandleFunc("/posts/count", api.countPosts).Methods("GET")
	router.HandleFunc("/posts/search", api.findPostsByCriteria).Methods("POST")
	router.HandleFunc("/posts/search/count", api.countPostsByCriteria).Methods("POST")
	router.HandleFunc("/posts/{id}", api.findPostByID).Methods("GET")

	return router
}
