package main

import (
	"github.com/lazylex/watch-store/secure/internal/repository/in_memory/redis"
	"github.com/lazylex/watch-store/secure/internal/repository/joint"
	"github.com/lazylex/watch-store/secure/internal/repository/persistent/postgresql"
	"github.com/lazylex/watch-store/secure/internal/service"
)

// runCodeForDebug код внутри функции запускается при разработке или отладке приложения. Данный файл будет помещен в
// gitignore и в дальнейшем не будет коммитится в репозиторий
func runCodeForDebug(memory *redis.Redis, persistent *postgresql.PostgreSQL, joint joint.Repository, serv *service.Service) {

}
