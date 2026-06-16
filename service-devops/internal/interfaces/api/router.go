package api

import (
	"github.com/gin-gonic/gin"
	"leap-one/service-devops/internal/interfaces/api/handler"
)

func RegisterRoutes(
	r *gin.Engine,
	repoH *handler.RepoHandler,
	pipeH *handler.PipelineHandler,
	artiH *handler.ArtiHandler,
	deployH *handler.DeployHandler,
	envH *handler.EnvHandler,
) {
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "leap-one-devops"})
	})

	v1 := r.Group("/api/v1")
	{
		// Repositories
		repos := v1.Group("/repositories")
		{
			repos.POST("", repoH.CreateRepo)
			repos.GET("", repoH.ListRepos)
			repos.GET("/:id", repoH.GetRepo)
			repos.PUT("/:id", repoH.UpdateRepo)
			repos.DELETE("/:id", repoH.DeleteRepo)
			repos.POST("/:id/test-connection", repoH.TestConnection)
		}

		// Pipelines
		pipes := v1.Group("/pipelines")
		{
			pipes.POST("", pipeH.CreatePipeline)
			pipes.GET("", pipeH.ListPipelines)
			pipes.GET("/:id", pipeH.GetPipeline)
			pipes.PUT("/:id", pipeH.UpdatePipeline)
			pipes.DELETE("/:id", pipeH.DeletePipeline)
			pipes.POST("/:id/trigger", pipeH.TriggerPipeline)
			pipes.GET("/:id/runs", pipeH.ListRuns)
			pipes.GET("/:id/runs/:rid", pipeH.GetRun)
			pipes.POST("/:id/runs/:rid/cancel", pipeH.CancelRun)
		}

		// Artifacts
		artis := v1.Group("/artifacts")
		{
			artis.GET("", artiH.ListArtifacts)
			artis.GET("/:id", artiH.GetArtifact)
			artis.DELETE("/:id", artiH.DeleteArtifact)
		}

		// Deployments
		deploys := v1.Group("/deployments")
		{
			deploys.POST("", deployH.CreateDeployment)
			deploys.GET("", deployH.ListDeployments)
			deploys.GET("/:id", deployH.GetDeployment)
			deploys.POST("/:id/rollback", deployH.RollbackDeployment)
		}

		// Environment Variables
		envs := v1.Group("/env-vars")
		{
			envs.GET("", envH.ListEnvVars)
			envs.POST("", envH.CreateEnvVar)
			envs.PUT("/:id", envH.UpdateEnvVar)
			envs.DELETE("/:id", envH.DeleteEnvVar)
		}
	}
}
