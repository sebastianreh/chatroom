package httpserver

// Routes build the routes of the server
func (s *Server) Routes() {
	root := s.Server.Group(s.dependencies.Config.Prefix)
	root.GET("/ping", s.dependencies.PingHandler.Ping)

	userGroup := root.Group("/user")
	userGroup.POST("", s.dependencies.UserHandler.Create)
	userGroup.POST("/login", s.dependencies.UserHandler.Login)
	userGroup.GET("", s.dependencies.UserHandler.Get)
	userGroup.DELETE("/:id", s.dependencies.UserHandler.Delete)

	roomGroup := root.Group("/room")
	roomGroup.POST("", s.dependencies.RoomHandler.Create)
	roomGroup.GET("", s.dependencies.RoomHandler.Get)
	roomGroup.DELETE("/:id", s.dependencies.RoomHandler.Delete)

	sessionGroup := root.Group("/session")
	sessionGroup.POST("/join", s.dependencies.SessionHandler.Join)
	sessionGroup.POST("/exit", s.dependencies.SessionHandler.Exit)
	sessionGroup.GET("/messages/:room_id", s.dependencies.SessionHandler.GetMessages)
	sessionGroup.GET("", s.dependencies.SessionHandler.HandleConnection)
}
