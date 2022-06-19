CREATE DATABASE IF NOT EXISTS `hexagony`; 

USE `hexagony`;

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `uuid` varchar(36) NOT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `users` WRITE;

UNLOCK TABLES;

LOCK TABLES `users` WRITE;

INSERT INTO `users` VALUES ('7d31461a-6ed5-425e-96fe-fa98e56d6828', 'John Doe', 'john@doe.com', '$2a$10$rPyJPskrTN545bXE0cqEU.T3uqluwiPFjGHMjE0/K.QuTe5XedjYi', '2022-06-19 16:53:09.000', '2022-06-19 16:53:09.000');

UNLOCK TABLES;

DROP TABLE IF EXISTS `albums`;

CREATE TABLE `albums` (
  `uuid` varchar(36) NOT NULL,
  `name` varchar(100) NOT NULL,
  `length` int(10) unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `albums` WRITE;

UNLOCK TABLES;
