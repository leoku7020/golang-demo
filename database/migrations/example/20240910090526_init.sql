-- +goose Up
-- 建立 members 表
CREATE TABLE members (
     id INT AUTO_INCREMENT PRIMARY KEY,
     username VARCHAR(255) NOT NULL,
     password VARCHAR(255) NOT NULL
);

-- 建立 items 表
CREATE TABLE items (
   id INT AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   category VARCHAR(255) NOT NULL
);

INSERT INTO members (username, password)
VALUES
    ('admin', '21232F297A57A5A743894A0E4A801FC3');

-- 插入 5 筆假資料到 items 表
INSERT INTO items (name, category)
VALUES
    ('Laptop', 'Electronics'),
    ('Chair', 'Furniture'),
    ('Book', 'Stationery'),
    ('Smartphone', 'Electronics'),
    ('Table', 'Furniture');

-- +goose Down
DROP TABLE members;
DROP TABLE items;
