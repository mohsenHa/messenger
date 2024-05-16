-- +migrate Up
-- please read this article to understand why we use VARCHAR(191)
-- https://www.grouparoo.com/blog/varchar-191#why-varchar-and-not-text
CREATE TABLE `users` (
                       `id` VARCHAR(191) PRIMARY KEY NOT NULL,
                       `public_key` TEXT NOT NULL,
                       `active_code` VARCHAR(50),
                       `status` INT NOT NULL DEFAULT 0
);

-- +migrate Down
DROP TABLE `users`;