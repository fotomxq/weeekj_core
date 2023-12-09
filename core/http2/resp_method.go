package CoreHttp2

// GET GET方法
func (t *Core) GET() *Core {
	t.Method = MethodGet
	t.Resp()
	t.DoResp()
	return t
}

// POST POST方法
func (t *Core) POST() *Core {
	t.Method = MethodPost
	t.Resp()
	t.DoResp()
	return t
}

// PUT PUT方法
func (t *Core) PUT() *Core {
	t.Method = MethodPut
	t.Resp()
	t.DoResp()
	return t
}

// PATCH PATCH方法
func (t *Core) PATCH() *Core {
	t.Method = MethodPatch
	t.Resp()
	t.DoResp()
	return t
}

// DELETE DELETE方法
func (t *Core) DELETE() *Core {
	t.Method = MethodDelete
	t.Resp()
	t.DoResp()
	return t
}
