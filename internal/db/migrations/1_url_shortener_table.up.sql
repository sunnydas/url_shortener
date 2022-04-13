CREATE TABLE `shortened_urls` (
                            `id` bigInt(12) NOT NULL AUTO_INCREMENT,
                            `original_url` MEDIUMTEXT NOT NULL,
                            `shortened_url` VARCHAR(26) NOT NULL,
                            `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `requester_id` VARCHAR(12) NOT NULL,
                            `expiry_date` DATETIME NOT NULL,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY (`shortened_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;