-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    username   VARCHAR(50)  NOT NULL UNIQUE,
    email      VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_username ON users (username);
CREATE INDEX idx_email ON users (email);

-- Добавление демо-данных
INSERT INTO users (username, email)
VALUES ('user1', 'user1@example.com'),
       ('user2', 'user2@example.com'),
       ('user3', 'user3@example.com'),
       ('john_doe', 'john@example.com'),
       ('jane_smith', 'jane@example.com'),
       ('tech_guru', 'tech@example.com'),
       ('marketing_pro', 'marketing@example.com'),
       ('design_master', 'design@example.com'),
       ('dev_ninja', 'developer@example.com'),
       ('data_scientist', 'data@example.com');

-- Создание пользователя с привилегиями (если нужны дополнительные привилегии)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON app_database.* TO 'app_user'@'%';
-- FLUSH PRIVILEGES;

########################################################################################################################
-- Миграция для создания таблицы products
-- Создаем таблицу для продуктов
CREATE TABLE IF NOT EXISTS products
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    quantity    INT            NOT NULL DEFAULT 0,
    created_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP               DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
);

-- Добавляем демонстрационные данные
INSERT INTO products (name, description, price, quantity, created_at, updated_at)
VALUES ('Смартфон Super Phone X5', 'Флагманский смартфон с 6.5" AMOLED экраном, 8 ГБ RAM и 128 ГБ ROM', 59999.99, 15,
        NOW(), NOW()),
       ('Ноутбук ProBook 15', '15.6" ноутбук для профессионалов с Intel Core i7, 16 ГБ RAM, SSD 512 ГБ', 89999.50, 8,
        NOW(), NOW()),
       ('Беспроводные наушники SoundBuds', 'Bluetooth наушники с активным шумоподавлением и 24 часами работы', 4999.90,
        25, NOW(), NOW()),
       ('Умные часы FitTrack Pro', 'Смарт-часы с мониторингом пульса, сна и GPS-трекером', 7499.00, 12, NOW(), NOW()),
       ('Планшет TabMax', '10.1" планшет с разрешением 2K, 6 ГБ RAM и батареей 8000 mAh', 24999.00, 10, NOW(), NOW()),
       ('Игровая мышь ProGamer', 'Эргономичная игровая мышь с RGB подсветкой и 7 программируемыми кнопками', 2999.50,
        30, NOW(), NOW()),
       ('Механическая клавиатура TypeMaster', 'Механическая клавиатура с синими переключателями и RGB подсветкой',
        6499.90, 18, NOW(), NOW()),
       ('Внешний жёсткий диск StorageMax 2TB', 'Портативный HDD на 2 ТБ с USB 3.1 и шифрованием данных', 5999.00, 22,
        NOW(), NOW()),
       ('Wi-Fi роутер NetMaster AC1200', 'Двухдиапазонный Wi-Fi роутер с технологией Mesh и гигабитными портами',
        3999.90, 14, NOW(), NOW()),
       ('Монитор ViewPro 27"', '27" IPS монитор с разрешением 4K, 144 Гц и поддержкой HDR', 29999.00, 7, NOW(), NOW()),
       ('Портативная колонка SoundBox', 'Беспроводная портативная колонка с защитой IPX7 и 20 часами работы', 3499.90,
        20, NOW(), NOW()),
       ('Графический планшет DrawTab', 'Планшет для художников с 8192 уровнями нажатия и беспроводным пером', 8499.00,
        9, NOW(), NOW()),
       ('Фитнес-браслет HealthBand', 'Водонепроницаемый фитнес-трекер с мониторингом пульса и сна', 1999.90, 35, NOW(),
        NOW()),
       ('Веб-камера StreamPro 4K', 'Веб-камера с разрешением 4K, автофокусом и шумоподавляющими микрофонами', 4499.00,
        16, NOW(), NOW()),
       ('Внешний аккумулятор PowerMax 20000 mAh',
        'Портативное зарядное устройство с функцией быстрой зарядки и 3 USB портами', 2799.50, 28, NOW(), NOW());

-- Процедура для обновления количества продуктов при заказе
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS update_product_quantity(
    IN product_id BIGINT,
    IN quantity_change INT
)
BEGIN
    UPDATE products
    SET quantity = quantity + quantity_change
    WHERE id = product_id
      AND (quantity + quantity_change) >= 0;

    SELECT ROW_COUNT() > 0 AS success;
END //
DELIMITER ;

-- Индекс для быстрого поиска по цене
CREATE INDEX idx_price ON products (price);

-- Представление для быстрого доступа к товарам в наличии
CREATE OR REPLACE VIEW products_in_stock AS
SELECT id, name, description, price, quantity
FROM products
WHERE quantity > 0
ORDER BY name;

-- Представление для быстрого доступа к популярным товарам (на основе цены, в будущем можно изменить на основе продаж)
CREATE OR REPLACE VIEW popular_products AS
SELECT id, name, description, price, quantity
FROM products
WHERE quantity > 0
ORDER BY price DESC
LIMIT 10;

########################################################################################################################

CREATE TABLE purchases (
                           id BIGINT AUTO_INCREMENT PRIMARY KEY,
                           user_id BIGINT NOT NULL,
                           product_id BIGINT NOT NULL,
                           quantity INT NOT NULL,
                           total_price DECIMAL(10, 2) NOT NULL,
                           status VARCHAR(20) NOT NULL,
                           created_at TIMESTAMP NOT NULL,
                           updated_at TIMESTAMP NOT NULL,
                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                           FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
                           INDEX (user_id),
                           INDEX (product_id),
                           INDEX (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;