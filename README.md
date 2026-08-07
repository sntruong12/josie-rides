# Josie Rides

This repository follows the book Let's Go by Alex Edwards. It houses golang code to be deployed on AWS. I'm using this as practice for modern cloud development. 

## Golang Notes

| CLI Command | Description |
| ----------- | ----------- |
| go mod init nameOfPackage | creates your golang project as a module, creates a document with all the project dependencies, your module's path, and module's go version |
| go run pathToPackage | a shortcut that compiles your code, creates an executable binary in your /tmp directory, and then runs this binary in one step |

## homebrew notes

| CLI command | description |
| ----------- | ----------- |
| brew services list | lists all the managed services and their status |
| brew services start mysql | starts the mysql server |
| brew services stop mysql | stops the mysql server |

## MySQL notes

| SQL query/CLI command | description |
| ----------- | ----------- |
| CREATE DATABASE databaseName CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; | Creates a new UTF-8 database |
| USE databaseName; | Switches to the database named databaseName. |
| SHOW DATABASES; | Lists all available databases on the MySQL server. |
| SHOW TABLES; | Lists all tables in the current database. |
| CREATE USER 'userName'@'localhost' IDENTIFIED BY 'password'; | Creates a new user named userName that can only login from localhost with the password password |
| GRANT SELECT, INSERT, UPDATE, DELETE ON databaseName.* TO 'userName'@'localhost'; | Grants specific privileges on databaseName to userName |
| mysql -D databaseName -u userName -p | login to the database named databaseName as userName |

## Shell Commands

| Command | Description |
| ------- | ----------- |
| curl -i -X POST url | sends and http request with the POST method and include the header in the output |
| curl | in general wrap the url parameter in single quotes so that you won't have to escape some special characters |
| jobs | displays a list of currently executing jobs/processes |
| fg | get the job that is currently running in background to the foreground, allowing you to interact with the app in the terminal |

## Project Structure

| directory/file | description |
| -------------- | ----------- |
| cmd | The cmd directory will contain the application-specific code for the executable applications in the project. For now we’ll have just one executable application — the web application — which will live under the cmd/web directory. |
| internal | The internal directory will contain the ancillary non-application-specific code used in the project. We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project. |
| ui | The ui directory will contain the user-interface assets used by the web application. Specifically, the ui/html directory will contain HTML templates, and the ui/static directory will contain static files (like CSS and images).|

Two benefits for using this structure
1. Clean separation of Go and non-Go assets. Go code lives in cmd and internal directories, leaving the root directory free to hold non-Go assets like UI files, makefiles, module definitions like our go.mod file. Easier to manage when it comes to building and deploying the app.
2. Scales like if you want to add more executable application to the project. Adding a CLI to automate some admin tasks, would be able to reuse internal code keep your code DRY.

Some tradeoffs here
1. Blast Radius, if multiple services share the internal code, a bug in internal could affect multiple services like the web server and a cloud function all at once.
2. Slower CI/CD pipelines, new updates to code would have to be handled specifically. CI flows would have to figure what changed and which services need to be built. Tests and builds can take longer, thus impacting speed of delivery.

## Data Models
1. table name, rides
id
title
description
trail name
distance(imperial)
created at (time when blog post written)
date time of ride (rode at)
duration (of the ride)
media(json object with videos, images and whatever else we decide to store - where value for each key is an array of strings representing the url for where the image is hosted, GCS/S3)

```
-- Create a new UTF-8 `josie_rides` database.
CREATE DATABASE josie_rides CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- switch to the josie_rides db
USE josie_rides;

-- create table query
CREATE TABLE rides (
    id                  INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title               VARCHAR(255) NOT NULL,
    description         LONGTEXT NOT NULL,
    trail_name          VARCHAR(255) NULL,
    distance_miles      DECIMAL(5,2) NOT NULL,
    duration            TIME NOT NULL,
    rode_at             DATETIME NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    media               JSON NULL,

    INDEX idx_rode_at (rode_at),
    INDEX idx_trail_name (trail_name),

    CONSTRAINT chk_distance_miles_positive CHECK (distance_miles > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- list of indexes

-- implemented

-- Index for filtering or ordering by ride date/time
CREATE INDEX idx_rode_at ON rides (rode_at);
-- Index for searching or filtering rides by trail name
CREATE INDEX idx_trail_name ON rides (trail_name);

-- not implemented

-- Index for sorting blog posts chronologically by creation date
CREATE INDEX idx_created_at ON rides (created_at);
-- Full-text index for searching post titles and descriptions
CREATE FULLTEXT INDEX idx_ft_search ON rides (title, description);

-- seed database with dummy data
INSERT INTO rides (
    title, 
    description, 
    trail_name, 
    distance_miles, 
    duration, 
    rode_at, 
    media
) VALUES 
(
    'Morning Loop at Alviso Marina', 
    'A scenic and flat gravel ride around the salt ponds. Great weather with light headwinds on the return leg.', 
    'Alviso Marina Loop', 
    14.50, 
    '01:12:30', 
    '2026-07-15 08:30:00', 
    JSON_OBJECT(
        'images', JSON_ARRAY('https://storage.googleapis.com/josie-rides-assets/alviso1.jpg', 'https://storage.googleapis.com/josie-rides-assets/alviso2.jpg'),
        'videos', JSON_ARRAY('https://storage.googleapis.com/josie-rides-assets/alviso_clip.mp4')
    )
),
(
    'Challenging Climb up Mount Hamilton', 
    'Tough climb with rewarding views of the Bay Area. Pushed hard on the final switchbacks.', 
    'Mount Hamilton Road', 
    38.25, 
    '03:45:15', 
    '2026-07-22 07:00:00', 
    JSON_OBJECT(
        'images', JSON_ARRAY('https://storage.googleapis.com/josie-rides-assets/mthamilton_summit.jpg')
    )
),
(
    'Sunset Coastal Cruise', 
    'Easy recovery ride along the coastline. Stopped for coffee midway.', 
    'Half Moon Bay Coastal Trail', 
    9.75, 
    '00:48:10', 
    '2026-08-02 18:15:00', 
    JSON_OBJECT(
        'images', JSON_ARRAY('https://storage.googleapis.com/josie-rides-assets/hmb_sunset.jpg')
    )
);

CREATE USER 'web'@'localhost' IDENTIFIED BY 'password';
GRANT SELECT, INSERT, UPDATE, DELETE ON josie_rides.* TO 'web'@'localhost';

```

## General Notes
1. Continuous Integration is testing the integrity of the code when code changes happen. Typically test then build the code if it's tested fine.
2. Continuous Delivery is releasing the tested code into production. 
