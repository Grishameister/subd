create database forum
	with owner docker
	encoding 'utf8'
    TABLESPACE = pg_default
	;
GRANT ALL PRIVILEGES ON database forum TO docker;