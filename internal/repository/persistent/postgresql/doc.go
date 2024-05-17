/*
Package postgresql: пакет для осуществления взаимодействия с СУБД PostgreSQL. Общение с БД осуществляется через пул
соединений, доступный посредством методов из пакета 'github.com/jackc/pgx'. Методы для взаимодействия с БД содержит
структура PostgreSQL. Функция MustCreate возвращает заполненную структуру PostgreSQL в случае успешной установки связи с
базой данных. В противном случае выполнение приложения прекращается. При отсутствии в базе данных схемы или какой-либо
из необходимых для работы таблиц, они создаются.
*/
package postgresql
