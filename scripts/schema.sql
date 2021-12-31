CREATE DATABASE IF NOT EXISTS app;
use app;  


CREATE TABLE IF NOT EXISTS employees (
uid varchar(50) PRIMARY KEY NOT NULL,
name VARCHAR(30) NOT NULL,
gender INT NOT NULL,
dob VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
uid varchar(50) PRIMARY KEY NOT NULL,
name VARCHAR(30) NOT NULL,
description VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS members(
employee_id varchar(50) NOT NULL,
team_id varchar(50) NOT NULL,
PRIMARY KEY (employee_id, team_id)
);