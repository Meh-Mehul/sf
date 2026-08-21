package handlers

type SF_SessionContext struct {
	Id          string
	IsTerm      bool
	TTL         int
	LastUpdated int
}

type SF_Session struct {
	Context SF_SessionContext
}
