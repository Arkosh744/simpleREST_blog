package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Posts interface {
	Create(ctx context.Context, post domain.Post) error
	GetById(ctx context.Context, id int64) (domain.Post, error)
	GetAll(ctx context.Context) ([]domain.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, post *domain.UpdatePost) error
}

type Handler struct {
	postServices Posts
}

func NewHandler(posts Posts) *Handler {
	return &Handler{
		postServices: posts,
	}
}

func (h *Handler) InitRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/post/new", loggerMiddleware(h.NewPost))
	mux.HandleFunc("/post/all", loggerMiddleware(h.GetAllPosts))
	mux.HandleFunc("/post/get/", loggerMiddleware(h.GetPostById))
	mux.HandleFunc("/post/update", loggerMiddleware(h.UpdatePostById))
	mux.HandleFunc("/post/delete", loggerMiddleware(h.DeletePostById))
	return mux
}

func (h *Handler) NewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var post domain.Post
		if err = json.Unmarshal(reqBytes, &post); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.postServices.Create(context.TODO(), post)
		if err != nil {
			log.Println("create() error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "application/json")
		posts, err := h.postServices.GetAll(context.TODO())
		if err != nil {
			log.Println("getAllPosts() error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(posts)
		if err != nil {
			log.Println("getAllPosts() error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			errorresponse := domain.PostError{"Error getting post: " + err.Error()}
			response, _ := json.Marshal(errorresponse)
			w.Write(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "application/json")

		id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/post/get/"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorresponse := domain.PostError{"Invalid id - ensure it is a number"}
			response, _ := json.Marshal(errorresponse)
			w.Write(response)
			return
		}
		posts, err := h.postServices.GetById(context.TODO(), id)
		if err != nil {
			if errors.Is(err, domain.ErrPostNotFound) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("getPostbyId() error:", err)
			errorresponse := domain.PostError{"Error getting post: " + err.Error()}
			response, _ := json.Marshal(errorresponse)
			w.Write(response)
			return
		}

		response, err := json.Marshal(posts)
		if err != nil {
			log.Println("getPostbyId() error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			errorresponse := domain.PostError{"Error getting post: " + err.Error()}
			response, _ := json.Marshal(errorresponse)
			w.Write(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) UpdatePostById(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var post *domain.UpdatePost
		if err = json.Unmarshal(reqBytes, &post); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.postServices.Update(context.TODO(), post.Id, post)
		if err != nil {
			log.Println("update() error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeletePostById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var post *domain.UpdatePost
		if err = json.Unmarshal(reqBytes, &post); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.postServices.Delete(context.TODO(), post.Id)
		if err != nil {
			log.Println("delete() error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
