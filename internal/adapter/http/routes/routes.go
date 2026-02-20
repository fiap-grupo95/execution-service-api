package routes

import (
	"strconv"

	_ "github.com/fiap-grupo95/execution-service-api/docs"

	"github.com/fiap-grupo95/execution-service-api/internal/adapter/gateway"
	handlers "github.com/fiap-grupo95/execution-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/execution-service-api/internal/adapter/http/middleware"
	execution_repository "github.com/fiap-grupo95/execution-service-api/internal/adapter/persistence/mongodb"
	"github.com/fiap-grupo95/execution-service-api/internal/infrastructure/database"
	"github.com/fiap-grupo95/execution-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/execution-service-api/internal/infrastructure/observability"
	"github.com/fiap-grupo95/execution-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

const PORT = 8081

func Run() {
	newRelicApp, err := observability.NewRelicApp()
	if err != nil {
		logs.Logger().Error().Err(err).Msg("Failed to initialize New Relic; continuing without New Relic integration")
	}

	setMiddlewares()

	router.Use(nrgin.Middleware(newRelicApp))

	logs.Init(newRelicApp)

	if newRelicApp != nil {
		observability.SetMetricsCollector(
			observability.NewNewRelicMetricsCollector(newRelicApp),
		)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	InitApp()

	logger := logs.Logger()
	err = router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func InitApp() {
	db, err := database.ConnectDatabase()
	if err != nil {
		logs.Logger().Fatal().Err(err).Msg("failed to connect to MongoDB")
	}

	executionServiceRepository := execution_repository.NewExecutionServiceRepository(db)

	executionGateway := gateway.NewExecutionServiceGateway(executionServiceRepository)
	executionUsecase := usecase.NewExecutionUsecase(executionGateway)

	executionHandler := handlers.NewExecutionHandler(executionUsecase)

	addPingRoutes(router)
	addExecutionRoutes(router, executionHandler)
}

func setMiddlewares() {

	middleware.SetTrustedProxies(router)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logs.Logger().Error().Msgf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))

}
