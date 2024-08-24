package rest

import (
	"github.com/SavelyDev/crud-app/internal/domain"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/SavelyDev/crud-app/docs"
)

type Auth interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accesToken string) (int, error)
}

type TodoList interface {
	CreateList(userId int, todoList domain.TodoList) (int, error)
	GetAllLists(userId int) ([]domain.TodoList, error)
	GetListById(userId, listId int) (domain.TodoList, error)
	DeleteList(userId, listId int) error
	UpdateList(userId, listId int, input domain.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(userId, listId int, input domain.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]domain.TodoItem, error)
	GetItemById(userId, itemId int) (domain.TodoItem, error)
	DeleteItem(userId, itemId int) error
	UpdateItem(userId, itemId int, input domain.UpdateItemInput) error
}

type Handler struct {
	AuthService     Auth
	TodoListService TodoList
	TodoItemService TodoItem
}

func NewHandler(auth Auth, todoList TodoList, todoItem TodoItem) *Handler {
	return &Handler{AuthService: auth,
		TodoListService: todoList,
		TodoItemService: todoItem,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}
