ALTER DATABASE valhaladb CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE valhaladb.employees (
    id MEDIUMINT NOT NULL AUTO_INCREMENT,
    name varchar(40),
    PRIMARY KEY (id)
);

CREATE TABLE valhaladb.client (
    id MEDIUMINT NOT NULL AUTO_INCREMENT,
    name varchar(40),
    email varchar(40),
    phone varchar(10),
    PRIMARY KEY (id)
);




SET names utf8;
INSERT INTO valhaladb.employees(id, name) VALUES (1, "Schelb");
INSERT INTO valhaladb.employees(id, name) VALUES (2, "Tchelão");
INSERT INTO valhaladb.employees(id, name) VALUES (3, "Dudu");