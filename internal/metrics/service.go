package metrics

import "github.com/prometheus/client_golang/prometheus"

type Service struct {
	login               *prometheus.CounterVec
	logout              *prometheus.CounterVec
	authenticationError *prometheus.CounterVec
}

func (s *Service) AuthenticationErrorInc() {
	s.authenticationError.With(prometheus.Labels{}).Inc()
}

func (s *Service) LoginInc() {
	s.login.With(prometheus.Labels{}).Inc()
}

func (s *Service) LogoutInc() {
	s.logout.With(prometheus.Labels{}).Inc()
}

// createLoginTotalMetric создает и регистрирует метрику login_total, являющуюся счетчиком залогиненых пользователей
// (сервисов)
func createLoginTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	login := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "login_total",
		Namespace: NAMESPACE,
		Help:      "Count of login users (services)",
	}, []string{})
	if err = prometheus.Register(login); err != nil {
		return nil, err
	}

	login.With(prometheus.Labels{})

	return login, nil
}

// createLogoutTotalMetric создает и регистрирует метрику logout_total, являющуюся счетчиком вышедших из сеанса (не по
// таймауту) пользователей (сервисов)
func createLogoutTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	logout := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "logout_total",
		Namespace: NAMESPACE,
		Help:      "Count of logout users (services)",
	}, []string{})
	if err = prometheus.Register(logout); err != nil {
		return nil, err
	}

	logout.With(prometheus.Labels{})

	return logout, nil
}

// createAuthenticationErrorTotalMetric создает и регистрирует метрику authentication_error_total, являющуюся счетчиком
// ошибок логи́на пользователей (сервисов)
func createAuthenticationErrorTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	authErrors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "authentication_error_total",
		Namespace: NAMESPACE,
		Help:      "Count of authentication errors",
	}, []string{})
	if err = prometheus.Register(authErrors); err != nil {
		return nil, err
	}

	authErrors.With(prometheus.Labels{})

	return authErrors, nil
}
