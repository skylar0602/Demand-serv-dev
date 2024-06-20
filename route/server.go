package route

import (
	"fmt"
	"net/http"

	"github.com/smarterwallet/demand-abstraction-serv/pkg"

	log "github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smarterwallet/demand-abstraction-serv/config"
	"github.com/smarterwallet/demand-abstraction-serv/model"
	"github.com/smarterwallet/demand-abstraction-serv/service"
)

type HTTPServer struct {
	*gin.Engine
	config    *config.Config
	demandSrv *service.DemandService
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	engine := gin.Default()
	pkg.NewBaseService(cfg)
	srv := service.NewDemandService(cfg)
	if srv == nil {
		panic("NewDemandService error")
	}
	return &HTTPServer{
		Engine:    engine,
		config:    cfg,
		demandSrv: srv,
	}
}

func (s *HTTPServer) Start() {
	listenAddr := fmt.Sprintf(":%d", s.config.Port)
	// v1 := s.Group("/v1", CORSMiddleware())
	v1 := s.Group("/v1")
	v1.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})
	v1.POST("/ctx", s.ConversationHistory(), func(ctx *gin.Context) {
		cid, ok := ctx.Get(model.ConversationID)
		if !ok {
			SendErrorResponse(ctx, http.StatusInternalServerError, errors.New("missing cid"))
			return
		}
		var request model.CtxRequest
		if err := ctx.BindJSON(&request); err != nil {
			SendErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		err := s.demandSrv.InitCtx(ctx, cid.(string), &request)
		if err != nil {
			SendErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header(model.CIDHeader, cid.(string))
		ctx.JSON(200, map[string]interface{}{
			"code":    200,
			"message": "success",
			"result":  "ok",
		})
	})
	v1.POST("/chat", s.ConversationHistory(), func(ctx *gin.Context) {
		cid, ok := ctx.Get(model.ConversationID)
		if !ok {
			SendErrorResponse(ctx, http.StatusInternalServerError, errors.New("missing cid"))
			return
		}
		var request model.DemandRequest
		if err := ctx.BindJSON(&request); err != nil {
			SendErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		if request.Model != model.ModelV1 {
			SendErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid model"))
			return
		}
		resp, err := s.demandSrv.ChatDemand(ctx, cid.(string), request.Demand)
		if err != nil {
			SendErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Header(model.CIDHeader, cid.(string))
		ctx.JSON(200, resp)
	})
	log.Infof("server listen on: %s", listenAddr)
	go func() {
		if err := s.Run(listenAddr); err != nil && err != http.ErrServerClosed {
			log.Infof("listen: %s", err)
		}
	}()
}

func (s *HTTPServer) Stop() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// todo release resource
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *HTTPServer) ConversationHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.Request.Header.Get(model.CIDHeader)
		if cid == "" || cid == "null" || cid == "undefined" {
			cid = s.demandSrv.NewChat()
			c.Header(model.CIDHeader, cid)
			c.JSON(200, map[string]string{
				"cid": cid,
			})
			c.Abort()
			return
		}
		c.Set(model.ConversationID, cid)
		c.Next()
	}
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func SendErrorResponse(c *gin.Context, code int, err error) {
	if err == nil {
		c.AbortWithStatus(code)
		return
	}
	response := &ErrorResponse{}
	response.Error.Code = code
	response.Error.Message = err.Error()
	c.JSON(code, response)
	c.Abort()
}
