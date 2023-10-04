USE mysql;

DROP DATABASE IF EXISTS go_mysql_connection;
CREATE DATABASE go_mysql_connection;

DROP USER IF EXISTS 'gouser'@'%';
CREATE USER IF NOT EXISTS 'gouser'@'%';
GRANT ALL PRIVILEGES ON go_mysql_connection.* TO 'gouser'@'%';
ALTER USER 'gouser'@'%' IDENTIFIED WITH caching_sha2_password BY 'secret';
FLUSH PRIVILEGES;

USE go_mysql_connection;

CREATE TABLE products (
    id MEDIUMINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    price FLOAT(3,2) DEFAULT 0.00,
    sku CHAR(11) NOT NULL,
    updated_at DATETIME DEFAULT now() NOT NULL ON UPDATE now(),
    PRIMARY KEY (id)
) ENGINE=INNODB CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE components (
    id BIGINT NOT NULL AUTO_INCREMENT,
    product_id MEDIUMINT,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON UPDATE SET NULL
        ON DELETE CASCADE
) ENGINE=INNODB CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

DELIMITER |
CREATE TRIGGER sku_check BEFORE INSERT ON products
    FOR EACH ROW
BEGIN
    IF (NEW.sku REGEXP '^([0-9]{3})+-([0-9]{3})+-([0-9]{3})$') = 0 THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'SKU has wrong format!';
    END IF;
END;
|
DELIMITER ;

INSERT INTO products (id, name, price, sku, updated_at)
VALUES
    (1, 'Latte', 1.49, '123-456-789', now()),
    (2, 'Espresso', 0.99, '000-000-001', now());

INSERT INTO components (id, product_id, name)
VALUES
    (1, 1, 'Water'),
    (2, 1, 'Milk'),
    (3, 2, 'Chocolate');
