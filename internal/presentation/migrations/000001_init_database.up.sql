BEGIN;

-- Создание таблицы Users
CREATE TABLE Users (
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(256) UNIQUE NOT NULL,
    password    VARCHAR(128) NOT NULL,
    coins       INT NOT NULL
);

-- Создание индекса на email в таблице User
CREATE INDEX email_idx ON Users(email);

-- Создание таблицы Subject
CREATE TABLE Subject (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(32) NOT NULL,
    cost    INT NOT NULL
);

-- Создание индекса на name в таблице Subject
CREATE INDEX name_idx ON Subject(name);

-- Создание таблицы Inventory
CREATE TABLE Inventory (
    id           SERIAL PRIMARY KEY,
    subject_name VARCHAR(128) NOT NULL,
    user_id      INT NOT NULL,
    FOREIGN KEY (subject_name) REFERENCES Subject(id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Создание индекса на subject_name в таблице Inventory
CREATE INDEX subject_name_idx ON Inventory(subject_name);
-- Создание индекса на user_id в таблице Inventory
CREATE INDEX user_id_idx ON Inventory(user_id);

-- Создание таблицы Transaction
CREATE TABLE Transaction (
    id            SERIAL PRIMARY KEY,
    sender_name   VARCHAR(256) NOT NULL,
    receiver_name VARCHAR(256) NOT NULL,
    amount        INT NOT NULL,
    FOREIGN KEY (sender_name) REFERENCES Users(id),
    FOREIGN KEY (receiver_name) REFERENCES Users(id)
);

-- Создание индекса на sender_name в таблице Transaction
CREATE INDEX sender_name_idx ON Transaction(sender_name);
-- Создание индекса на receiver_name в таблице Transaction
CREATE INDEX receiver_name_idx ON Transaction(receiver_name);

-- Создание таблицы для хранения токенов доступа
CREATE TABLE Token (
    id      SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    value   VARCHAR(256) NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Создание индекса на value в таблице Token
CREATE INDEX value_idx ON Token(value);

-- Заполнение таблицы Subject товарами
INSERT INTO Subject (name, cost) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);

-- Коммит транзакции, если нет ошибок
COMMIT;

-- Роллбэк транзакции, если произошли ошибки
ROLLBACK;