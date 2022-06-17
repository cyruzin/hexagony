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
