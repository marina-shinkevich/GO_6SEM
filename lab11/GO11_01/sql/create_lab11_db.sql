IF DB_ID(N'PIS_Lab5') IS NULL
BEGIN
    CREATE DATABASE [PIS_Lab5];
END
GO

USE [PIS_Lab5];
GO

IF OBJECT_ID(N'dbo.celebrities', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.celebrities (
        id INT NOT NULL PRIMARY KEY,
        full_name NVARCHAR(200) NOT NULL,
        nationality NVARCHAR(120) NOT NULL,
        req_photo_path NVARCHAR(260) NOT NULL
    );
END
GO

IF NOT EXISTS (SELECT 1 FROM dbo.celebrities)
BEGIN
    INSERT INTO dbo.celebrities (id, full_name, nationality, req_photo_path)
    VALUES
        (1, N'Charlie Chaplin', N'United Kingdom', N'photos/charlie_chaplin.jpg'),
        (2, N'Audrey Hepburn', N'Belgium', N'photos/audrey_hepburn.jpg'),
        (3, N'Bruce Lee', N'China', N'photos/bruce_lee.jpg');
END
GO
