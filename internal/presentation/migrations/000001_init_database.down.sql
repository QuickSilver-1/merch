BEGIN;

-- Удаление таблиц, если они существуют
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Subject;
DROP TABLE IF EXISTS Inventory;
DROP TABLE IF EXISTS Transaction;
DROP TABLE IF EXISTS Token;

-- Коммит транзакции, если нет ошибок
COMMIT;

-- Роллбэк транзакции, если произошли ошибки
ROLLBACK;
