    CREATE DATABASE Auth_Dev;
    use Auth_Dev;
CREATE TABLE Role (
	RoleId int NOT NULL AUTO_INCREMENT primary Key,
    Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
    Is_Visible int NOT NULL DEFAULT(1), 
	RoleName varchar(100) NOT NULL
);

CREATE TABLE RestrictionEntity (
	EntityId int NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
    Is_Visible int NOT NULL DEFAULT(1), 
    EntityName varchar(100) NOT NULL
);

CREATE TABLE RestrictionEntityType (
	AuthTypeId int NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
    Is_Visible int NOT NULL DEFAULT(1), 
    AuthName varchar(100) NOT NULL
);

create Table User (
	UserId int NOT NULL AUTO_INCREMENT primary Key,
	Name varchar(100) NOT NULL,
    Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
    Is_Visible int NOT NULL DEFAULT(1), 
    email varchar(100) NOT NULL
);


Create Table UserCred (
	CredId int NOT NULL AUTO_INCREMENT primary Key,
	UserId int NOT NULL ,
    Hash_Password varchar(100) NOT NULL ,
    Is_Visible int NOT NULL DEFAULT(1), 
    Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
	FOREIGN KEY (UserId) REFERENCES User(UserId)
) ; 

Create Table RestrictionPerUser (
	RestrictId int NOT NULL AUTO_INCREMENT primary Key,
    UserId int NOT NULL ,
    EntityId int NOT NULL , 
    AuthTypeId int NOT NULL,
    Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
    Is_Visible int NOT NULL DEFAULT(1), 
    FOREIGN KEY (UserId) REFERENCES User(UserId),
    FOREIGN KEY (EntityId) REFERENCES RestrictionEntity(EntityId),
    FOREIGN KEY (AuthTypeId) REFERENCES RestrictionEntityType(AuthTypeId)
);

Create Table TokenStore (	
	UserId int NOT NULL primary Key,
    Token varchar(100) NOT NULL ,
	Last_Modified_Date DateTime NOT NULL DEFAULT (CURRENT_DATE()),
    Is_Visible int NOT NULL DEFAULT(1), 
    FOREIGN KEY (UserId) REFERENCES User(UserId)
) ;
    