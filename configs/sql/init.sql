create database forum
	with owner postgres
	encoding 'utf8'
    TABLESPACE = pg_default
	;
GRANT ALL PRIVILEGES ON database forum TO postgres;