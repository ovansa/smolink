package test

func (app *TestApp) SCleanup() {
	if app.pool != nil {
		if app.pgResource != nil {
			_ = app.pool.Purge(app.pgResource)
		}
		if app.redisResource != nil {
			_ = app.pool.Purge(app.redisResource)
		}
	}
}
