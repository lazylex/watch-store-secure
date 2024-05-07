package postgresql

// FillTestData заполняет таблицы БД данными для тестирования и отладки взаимодействия с БД не используя mock
func (p *PostgreSQL) FillTestData() error {
	stmt := `
	INSERT INTO services (service_id, name, description)
	VALUES (1, 'store', 'Оффлайн магазин'),
		   (2, 'stock', 'Склад'),
		   (3, 'secure', 'Сервис выдачи прав'),
		   (4, 'orders', 'Сервис обработки заказов'),
		   (5, 'map', 'Сервис поиска оптимальных маршрутов доставки'),
		   (6, 'scheduler', 'Планировщик закупок');
	
	INSERT INTO instances (instance_id, name, secret, service_fk)
	VALUES (1, 'Магазин в Донецке', 's1', 1),
		   (2, 'Магазин в Луганске','s2', 1),
		   (3, 'Магазин в Ростове', 's3',1),
		   (4, 'Магазин в Макеевке', 's4',1),
		   (5, 'Магазин в Матвеев-Кургане', 's5', 1),
		   (6, 'Склад в Донецке', 's6', 2),
		   (7, 'Склад в Ростове', 's7',2),
		   (8, 'Сервис выдачи прав', 's8', 3),
		   (9, 'Сервис поиска оптимальных маршрутов доставки', 's9', 5),
		   (10, 'Планировщик закупок', 's0', 6);
	
	INSERT INTO permissions (permission_id, name, number, description, service_fk)
	VALUES (1, 'убирать товар', 1, '', 1),
		   (2, 'добавлять товар', 2, '', 1),
		   (3, 'менять цену на товар', 3, '', 1),
		   (4, 'узнавать цену товара', 4, '', 1),
		   (5, 'узнавать количество товара', 5, '', 1),
		   (6, 'откладывать товар для заказа', 6, 'убирать товар из доступных к продаже, чтобы отправить заказчику', 1),
		   (7, 'продавать товар', 7, '', 1),
		   (8, 'оформлять возврат', 8, '', 1),
		   (9, 'узнавать количество товара', 1, '', 2),
		   (10, 'откладывать товар для заказа', 2, '', 2),
		   (11, 'откладывать товар для магазина', 3, '', 2),
		   (12, 'добавлять товар', 4, '', 2),
		   (13, 'убирать товар', 5, '', 2);
	
	INSERT INTO roles (role_id, name, description, service_fk)
	VALUES (1, 'Продавец', '', 1),
		   (2, 'Менеджер', '', 1),
		   (3, 'Сервис планирования', '', 1),
		   (4, 'Сервис заказа', '', 1),
		   (5, 'Сервис заказа', '', 2),
		   (6, 'Сервис планирования', '', 2);
	
	INSERT INTO groups (group_id, name, description, service_fk)
	VALUES (1, 'Персонал магазина', 'Продавцы и менеджеры', 1);
	
	INSERT INTO group_permissions (permission_fk, group_fk)
	VALUES (4, 1),
		   (5, 1);
	
	INSERT INTO role_permissions (role_fk, permission_fk)
	VALUES (1, 7),
		   (2, 1),
		   (2, 2);
	
	INSERT INTO group_roles (role_fk, group_fk)
	VALUES (1, 1),
		   (2, 1);
	
	INSERT INTO accounts (account_id, uuid, login, pwd_hash, state)
	VALUES (1, '0eca778b-d090-441a-bf29-be4f525f0b70', 'Камила Валиева',
			'$2a$14$IFAdmu0DP1PdOY3Is72zV.egTN.S9JaeuqRT0hadNvhTfwxuMRNbG', 1),
		   (2, '67249938-2743-442b-ad31-547c96274938', 'Алина Загитова',
			'$2a$14$xzS0lqHtX6vDD6zSK4Xo3eZIeqwdBXpOrDx34XERau9Ouhf0HxG8i', 1),
		   (3, '0acd9223-64de-4d21-8e18-2f3d61c9d5f7', 'Анна Щербакова',
			'$2a$14$ZkvAk8oY8iag.PVYHwHm4OR6K0BIG85dnPaEI8F0QOUXtBrE0hQIG', 1);
	
	INSERT INTO account_groups (account_fk, group_fk)
	VALUES (1, 1),
		   (2, 1),
		   (3, 1);
	
	INSERT INTO account_roles (role_fk, account_fk)
	VALUES (1, 1),
		   (1, 2),
		   (2, 3);
	
	INSERT INTO accounts_instances_permissions (account_fk, instance_fk, permission_fk)
	VALUES (3, 1, 8);`

	return p.processExecResult(p.pool.Exec(stmt))
}
