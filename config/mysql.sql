/* *****************************************************************************
// Setup the preferences
// ****************************************************************************/
SET NAMES utf8 COLLATE 'utf8_unicode_ci';
SET foreign_key_checks = 1;
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
SET default_storage_engine = InnoDB;
SET CHARACTER SET utf8;

/* *****************************************************************************
// Remove old database
// ****************************************************************************/
DROP DATABASE IF EXISTS Collector;

/* *****************************************************************************
// Create new database
// ****************************************************************************/
CREATE DATABASE Collector DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;
USE Collector;

/* *****************************************************************************
// Create the tables
// ****************************************************************************/
CREATE TABLE agency (

);

CREATE TABLE user (
    id INT(12) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_name VARCHAR(50) NOT NULL,
    password VARCHAR(60) NOT NULL,
    gender TINYINT(1) UNSIGNED DEFAULT 0,
    birth VARCHAR(50),
    mobile VARCHAR(50),
    agency_id INT(12) UNSIGNED NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT "customer",
    priority VARCHAR(20) NOT NULL DEFAULT "123",
    status TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_ontime DATETIME,
    last_onip VARCHAR(50) DEFAULT "",
    UNIQUE KEY (agency_id),
    PRIMARY KEY (id)
);

INSERT INTO `user` (`id`, `user_name`, `password`, `gender`, `birth`, `mobile`, `agency_id`, `role`, `priority`, `status`, `status`) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE note (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    content TEXT NOT NULL,
    
    user_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_note_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);