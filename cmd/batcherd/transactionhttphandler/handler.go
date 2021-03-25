package transactionhttphandler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/donut-batcher/internal/batch"
	"github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq"
	"github.com/mazxaxz/donut-batcher/pkg/message/transaction"
	"github.com/mazxaxz/donut-batcher/pkg/rest"
)

type handlerContext struct {
	batchSvc             batch.Service
	transactionPublisher rabbitmq.Publisher
	logger               *logrus.Logger
}

func New(bSvc batch.Service, transactionPublisher rabbitmq.Publisher, l *logrus.Logger) rest.SetupRouterer {
	c := handlerContext{
		batchSvc:             bSvc,
		transactionPublisher: transactionPublisher,
		logger:               l,
	}
	return &c
}

func (c *handlerContext) SetupRouter(r *gin.RouterGroup) {
	r.GET("/transactions/batches/history", c.GetBatchHistory)
	r.POST("/transactions/stress", c.Stress)
}

func (c *handlerContext) GetBatchHistory(cGin *gin.Context) {
	limit, err := strconv.Atoi(cGin.DefaultQuery("limit", "10"))
	if err != nil {
		httpErr := rest.NewError("invalid_parameter__limit", err)
		cGin.AbortWithStatusJSON(http.StatusBadRequest, httpErr)
		return
	}
	page, err := strconv.Atoi(cGin.DefaultQuery("page", "0"))
	if err != nil {
		httpErr := rest.NewError("invalid_parameter__page", err)
		cGin.AbortWithStatusJSON(http.StatusBadRequest, httpErr)
		return
	}
	order := cGin.DefaultQuery("order", "-1")
	status := cGin.Query("status")

	var asc bool
	if order == "-1" {
		asc = false
	} else {
		asc = true
	}

	var s *batch.Status
	if status != "" {
		v := batch.NewStatusFrom(status)
		s = &v
	}
	batches, err := c.batchSvc.Paginate(cGin, limit, limit*page, asc, s)
	if err != nil {
		httpErr := rest.NewError("paginate_error", err)
		cGin.AbortWithStatusJSON(http.StatusInternalServerError, httpErr)
		return
	}
	cGin.JSON(http.StatusOK, batches)
}

func (c *handlerContext) Stress(cGin *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cents := rand.Intn(max-min+1) + min
				msg := transaction.Transaction{
					ID:       uuid.NewString(),
					UserID:   fmt.Sprintf("user:%d", index),
					Amount:   fmt.Sprintf("%.2f", float64(cents)/float64(100)),
					Currency: "USD",
				}
				if err := c.transactionPublisher.Publish(cGin, msg, transaction.MessageTypeTransaction); err != nil {
					c.logger.Error(err)
				}
			}
		}(i)
	}
	wg.Wait()

	cGin.JSON(http.StatusNoContent, gin.H{})
}
