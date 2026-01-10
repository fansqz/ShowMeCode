package user

import (
	"github.com/fansqz/fancode-backend/controller/user"
	"github.com/gin-gonic/gin"
)

func SetupDebugRoutes(r *gin.Engine, debugController user.DebugController) {
	//用户相关
	judge := r.Group("/debug")
	{
		judge.POST("/session/create", debugController.CreateDebugSession)
		judge.GET("/sse/:id", debugController.CreateSseConnect)
		judge.POST("/start", debugController.Start)
		judge.POST("/step/in", debugController.StepIn)
		judge.POST("/step/out", debugController.StepOut)
		judge.POST("/step/over", debugController.StepOver)
		judge.POST("/continue", debugController.Continue)
		judge.POST("/sendToConsole", debugController.SendToConsole)
		judge.POST("/setBreakpoints", debugController.SetBreakpoints)
		judge.POST("/stackTrace", debugController.GetStackTrace)
		judge.POST("/frame/variables", debugController.GetFrameVariables)
		judge.POST("/variables", debugController.GetVariables)
		judge.POST("/terminate", debugController.Terminate)
	}
}
