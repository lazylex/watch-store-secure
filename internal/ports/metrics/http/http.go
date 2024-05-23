package http

//go:generate mockgen -source=http.go -destination=mocks/http.go
type MetricsInterface interface {
	RequestsTotalInc(map[string]string)
	RequestsDurationObserve(float64)
}
