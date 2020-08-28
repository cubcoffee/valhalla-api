ALTER DATABASE valhaladb CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE valhaladb.employees (
    id SMALLINT NOT NULL AUTO_INCREMENT,
    name varchar(40),
    responsibility varchar(20),
    hour_init TIME,
    hour_end TIME,
    credential_id INT,
    FOREIGN KEY (credential_id) REFERENCES valhaladb.credentials(id),
    PRIMARY KEY (id)
);

CREATE TABLE valhaladb.days_works (
    id SMALLINT NOT NULL AUTO_INCREMENT,
    day_index SMALLINT,
    user_id SMALLINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_employee_day
    FOREIGN KEY (user_id) REFERENCES valhaladb.employees(id)
);

SET names utf8;
INSERT INTO valhaladb.employees(id, name, responsibility, hour_init, hour_end)
VALUES 
    (1, "Schelb", "barbeiro","08:00:00", "18:00:00"),
    (2, "Tchelão", "barbeiro","08:00:00", "18:00:00"),
    (4, "Tchelão", "barbeiro","08:00:00", "18:00:00"),
    (3, "Dudu", "design","08:00:00", "18:00:00");

--  ODBC standard (1 = Sunday, 2 = Monday, ..., 7 = Saturday)
INSERT INTO valhaladb.days_works(day_index,user_id)
VALUES
    (1,1),
    (1,1),
    (2,1),
    (2,2),
    (4,3),
    (7,1);