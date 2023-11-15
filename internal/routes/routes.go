package routes

import (
	answerHandler "github.com/blazee5/quizmaster-backend/internal/answer/handler"
	answerRepo "github.com/blazee5/quizmaster-backend/internal/answer/repository"
	answerService "github.com/blazee5/quizmaster-backend/internal/answer/service"
	authHandler "github.com/blazee5/quizmaster-backend/internal/auth/handler"
	authRepo "github.com/blazee5/quizmaster-backend/internal/auth/repository"
	authService "github.com/blazee5/quizmaster-backend/internal/auth/service"
	questionHandler "github.com/blazee5/quizmaster-backend/internal/question/handler"
	questionRepo "github.com/blazee5/quizmaster-backend/internal/question/repository"
	questionService "github.com/blazee5/quizmaster-backend/internal/question/service"
	quizHandler "github.com/blazee5/quizmaster-backend/internal/quiz/handler"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz/repository"
	quizService "github.com/blazee5/quizmaster-backend/internal/quiz/service"
	userHandler "github.com/blazee5/quizmaster-backend/internal/user/handler"
	userRepo "github.com/blazee5/quizmaster-backend/internal/user/repository"
	userService "github.com/blazee5/quizmaster-backend/internal/user/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	log      *zap.SugaredLogger
	db       *sqlx.DB
	rdb      *redis.Client
	esclient *elasticsearch.Client
}

func NewServer(log *zap.SugaredLogger, db *sqlx.DB, rdb *redis.Client, esclient *elasticsearch.Client) *Server {
	return &Server{log: log, db: db, rdb: rdb, esclient: esclient}
}

func (s *Server) InitRoutes(e *echo.Echo) {
	authRepos := authRepo.NewRepository(s.db)
	authServices := authService.NewService(s.log, authRepos)
	authHandlers := authHandler.NewHandler(s.log, authServices)

	auth := e.Group("/auth")
	{
		auth.POST("/signup", authHandlers.SignUp)
		auth.POST("/signin", authHandlers.SignIn)
	}

	api := e.Group("/api")
	{
		userRepos := userRepo.NewRepository(s.db)
		userRedisRepo := userRepo.NewUserRedisRepo(s.rdb)
		userServices := userService.NewService(s.log, userRepos, userRedisRepo)
		userHandlers := userHandler.NewHandler(s.log, userServices)

		user := api.Group("/user", AuthMiddleware)
		{
			user.GET("/me", userHandlers.Get)
			user.GET("/:id", userHandlers.GetById)
			user.POST("/avatar", userHandlers.UploadAvatar)
			user.PUT("", userHandlers.Update)
			user.DELETE("", userHandlers.Delete)
		}

		quizRepos := quizRepo.NewRepository(s.db)
		quizRedisRepo := quizRepo.NewAuthRedisRepo(s.rdb)
		quizServices := quizService.NewService(s.log, quizRepos, quizRedisRepo, userRedisRepo)
		quizHandlers := quizHandler.NewHandler(s.log, quizServices)

		quiz := e.Group("/quiz")
		{
			quiz.POST("", quizHandlers.CreateQuiz, AuthMiddleware)
			quiz.GET("", quizHandlers.GetAllQuizzes)
			// quiz.GET("/search", quizHandlers.SearchByTitle)
			quiz.GET("/:id", quizHandlers.GetQuiz)
			quiz.POST("/:id/save", quizHandlers.SaveResult, AuthMiddleware)
			quiz.DELETE("/:id", quizHandlers.DeleteQuiz, AuthMiddleware)

			questionRepos := questionRepo.NewRepository(s.db)
			questionServices := questionService.NewService(s.log, questionRepos, quizRepos)
			questionHandlers := questionHandler.NewHandler(s.log, questionServices)

			question := quiz.Group("/:id/questions", AuthMiddleware)
			{
				question.POST("", questionHandlers.CreateQuestion)
				question.POST("/:questionId/image", questionHandlers.UploadImage)
				question.GET("", questionHandlers.GetQuizQuestions)
				question.PUT("/:questionId", questionHandlers.UpdateQuestion)
				question.DELETE("/:questionId", questionHandlers.DeleteQuestion)
				question.DELETE("/:questionId/image", questionHandlers.DeleteImage)

				answerRepos := answerRepo.NewRepository(s.db)
				answerServices := answerService.NewService(s.log, answerRepos, quizRepos)
				answerHandlers := answerHandler.NewHandler(s.log, answerServices)

				answer := question.Group("/:questionId/answers")
				{
					answer.POST("", answerHandlers.CreateAnswer)
					answer.PUT("/:answerId", answerHandlers.UpdateAnswer)
					answer.DELETE("/:answerId", answerHandlers.DeleteAnswer)
				}
			}
		}
	}

	e.Static("/public", "public")
}
