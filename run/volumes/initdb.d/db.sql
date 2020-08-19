ALTER DATABASE valhaladb CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE valhaladb.employees (
    id SMALLINT NOT NULL AUTO_INCREMENT,
    name varchar(40),
    responsibility varchar(20),
    PRIMARY KEY (id)
);

CREATE TABLE valhaladb.days_works (
    id SMALLINT NOT NULL AUTO_INCREMENT,
    day varchar(9),
    user_id SMALLINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_employee_day
    FOREIGN KEY (user_id) REFERENCES valhaladb.employees(id)
);

SET names utf8;
INSERT INTO valhaladb.employees(id, name, responsibility)
VALUES 
    (1, "Schelb", "barbeiro"),
    (2, "Tchelão", "barbeiro"),
    (3, "Dudu", "design");

INSERT INTO valhaladb.days_works(day,user_id)
VALUES
    ("Sunday", 1),
    ("Monday",1),
    ("Tuesday",1),
    ("Wednesday",2),
    ("Thursday",3),
    ("Saturday",1);