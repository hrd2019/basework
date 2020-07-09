package httpserver

import (
	"github.com/fuloge/basework/configs"
	"github.com/fuloge/basework/pkg/filter"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Server struct {
	engine     *gin.Engine
	pathGet    map[string]gin.HandlerFunc
	pathPost   map[string]gin.HandlerFunc
	pathPut    map[string]gin.HandlerFunc
	pathDelete map[string]gin.HandlerFunc
}

const (
	GET_METHOD    = 0
	POST_METHOD   = 1
	PUT_METHOD    = 2
	DELETE_METHOD = 3
)

func New() *Server {
	switch configs.EnvConfig.RunMode {
	case 1:
		gin.SetMode(gin.DebugMode)
	case 2:
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()
	e.Use(cross())

	filter := filter.Filter{}
	e.Use(filter.Checkauth())
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.NoRoute(noResponse)

	return &Server{
		engine:     e,
		pathPost:   make(map[string]gin.HandlerFunc),
		pathGet:    make(map[string]gin.HandlerFunc),
		pathPut:    make(map[string]gin.HandlerFunc),
		pathDelete: make(map[string]gin.HandlerFunc),
	}
}

func (s *Server) SetStatic(path string, dir string) {
	s.engine.Static(path, dir)
}

func (s *Server) LoadHTMLGlob(pattern string) {
	s.engine.LoadHTMLGlob(pattern)
}

func (s *Server) SetGetRouter(route string, handle func(*gin.Context)) {
	s.pathGet[route] = handle
}

func (s *Server) SetPostRouter(route string, handle func(*gin.Context)) {
	s.pathPost[route] = handle
}

func (s *Server) SetPutRouter(route string, handle func(*gin.Context)) {
	s.pathPut[route] = handle
}

func (s *Server) SetDeleteRouter(route string, handle func(*gin.Context)) {
	s.pathDelete[route] = handle
}

func (s *Server) assem() {
	authorized := s.engine.Group("/")
	{
		for key, handle := range s.pathGet {
			authorized.GET(key, handle)
		}

		for key, handle := range s.pathPost {
			authorized.POST(key, handle)
		}
	}
}

func (s *Server) SetGroup(route string) *gin.RouterGroup {
	return s.engine.Group(route)
}

func (s *Server) AddHandleByGroup(group *gin.RouterGroup, route string, methodType int, handle func(*gin.Context)) {
	switch methodType {
	case GET_METHOD:
		group.GET(route, handle)
	case POST_METHOD:
		group.POST(route, handle)
	case PUT_METHOD:
		group.PUT(route, handle)
	case DELETE_METHOD:
		group.DELETE(route, handle)
	default:
		print("no support method")
	}
}

func (s *Server) Run(ip string, port int64) {
	s.assem()
	s.engine.Run(ip + ":" + strconv.FormatInt(port, 10))
}

func cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		//c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			//c.AbortWithStatus(http.StatusNoContent)
			c.AbortWithStatus(http.StatusOK)
		}

		// 处理请求
		c.Next()
	}
}

func noResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "404, page not exists!",
	})
}
