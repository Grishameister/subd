create database forum
	with owner postgres
	encoding 'utf8'
	LC_COLLATE = 'ru_RU.UTF-8'
    LC_CTYPE = 'ru_RU.UTF-8'
    TABLESPACE = pg_default
	;
GRANT ALL PRIVILEGES ON database forum TO postgres;