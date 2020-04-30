CREATE TABLE sessions (
	secret text,
	dtmf text UNIQUE NOT NULL,
	purpose text,
	disclosed text,
	created timestamp NOT NULL DEFAULT now());
