package providers

//Provider interfaces all events providers
type Provider interface {
	GetAllEvents(config interface{}) ([]Event, error)
}
