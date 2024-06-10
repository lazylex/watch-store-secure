package postgresql

// FillTestData заполняет таблицы БД данными для тестирования и отладки взаимодействия с БД не используя mock.
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
	VALUES (1, 'получать данные о товаре', 1, '', 1),
		   (2, 'получать доступное для продажи количество товара', 2, '', 1),
		   (3, 'обновлять количество товара', 3, '', 1),
		   (4, 'обновлять цену товара', 4, '', 1),
		   (5, 'добавлять запись о товаре', 5, '', 1),
		   (6, 'получать общее количество проданного товара', 6, '', 1),
		   (7, 'осуществлять локальную продажу', 7, '', 1),
		   (8, 'резервировать товар', 8, '', 1),
		   (9, 'отменять резервирование', 9, '', 1),
		   (10, 'завершать продажу/отправку', 10, '', 1),

		   (11, 'узнавать количество товара', 1, '', 2),
		   (12, 'откладывать товар для заказа', 2, '', 2),
		   (13, 'откладывать товар для магазина', 3, '', 2),
		   (14, 'добавлять товар', 4, '', 2),
		   (15, 'убирать товар', 5, '', 2);
	
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
	VALUES (1, 1),
		   (2, 1);
	
	INSERT INTO role_permissions (role_fk, permission_fk)
	VALUES (1, 7),
	       (1, 10),
		   (2, 9),
		   (2, 10);
	
	INSERT INTO group_roles (role_fk, group_fk)
	VALUES (1, 1),
		   (2, 1);
-- 	Пароли Password_1, Password_2, Password_3:
	INSERT INTO accounts (account_id, uuid, login, pwd_hash, state)
	VALUES (1, '0eca778b-d090-441a-bf29-be4f525f0b70', 'Камила Валиева',
			'$2a$14$Ne17rB21.iXHWug6wuB80ethQ.vWrViWXpPFpUotkA8pkxAGqyAj2', 1),
		   (2, '67249938-2743-442b-ad31-547c96274938', 'Алина Загитова',
			'$2a$14$AIdOFdy1oMbExcRLRAq1Ju2cljtk96n20n4GaY14ZiYv0pFV1lOdS', 1),
		   (3, '0acd9223-64de-4d21-8e18-2f3d61c9d5f7', 'Анна Щербакова',
			'$2a$14$h0HKab7NHV8JHwCjUmwf9eVc.3k4d5jOoFutLudk1BwR2p3dYXTeS', 1),
	       (4, 'b05c8cb7-9d7b-4c12-a1c8-cc366514950c', 'store1',
			'$2a$14$DXeOD0RIEAQC6//0HLLqfe60EPILDblKQPNsmZE5NxurpkSq3bOQW', 1);
	
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
