 CREATE TABLE Credentials (
     id AUTO_INCREMENT PRIMARY KEY,
     password TEXT,
     salt TEXT,
 )

 CREATE TABLE IF NOT EXISTS credentials (
    id INT AUTO_INCREMENT,
    hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    PRIMARY KEY (id),
);