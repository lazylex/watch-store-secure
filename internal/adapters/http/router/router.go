package router

import "net/http"

var paths []string

// AssignPathToHandler проверяет, не прикреплен ли уже переданный первым аргументом функции адрес к какому-либо
// обработчику. Если прикреплен, то выполнение функции прекращается, чтобы не вызвать панику в http.HandleFunc. При
// нормальном выполнении, добавляет пусть к списку используемых и прикрепляет его к переданному третьим аргументом
// обработчику.
func AssignPathToHandler(path string, mux *http.ServeMux, handler func(http.ResponseWriter, *http.Request)) {
	for _, v := range paths {
		if v == path {
			return
		}
	}

	paths = append(paths, path)
	mux.HandleFunc(path, handler)
}

// ExistentPaths возвращает все зарегистрированные для обработчиков адреса.
func ExistentPaths() []string {
	return paths
}

// IsExistPath возвращает true, если в приложении используется передаваемый путь. Иначе - false.
func IsExistPath(path string) bool {
	for _, p := range ExistentPaths() {
		if p == path {
			return true
		}
	}

	return false
}
