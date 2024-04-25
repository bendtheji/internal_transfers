# Internal Transfers System

## Setup

### First, clone this repo.

```
git clone https://github.com/bendtheji/internal_transfers.git
```

### Install Go and MySQL if required
Install Go: https://go.dev/doc/install

Install MySQL: https://dev.mysql.com/downloads/installer/

### Setup MySQL Database
Enter MySQL shell as root:
```
mysql -u root -p
```
You'll need the password used during MySQL installation to log into the shell.

Create Database:
```
create database internal_transfers;
```

Create User and grant privileges for user to your created database:
```
CREATE USER <DB_USERNAME>@'localhost' identified by <DB_PASSWORD>;
GRANT ALL PRIVILEGES ON internal_transfers.* to <DB_USERNAME>@localhost with grant option;
```

Create tables for application:
```
CREATE TABLE `accounts` (
    `id` int NOT NULL AUTO_INCREMENT,
    `balance` decimal(10,2) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `transactions` (
    `source_account_id` int DEFAULT NULL,
    `destination_account_id` int DEFAULT NULL,
    `transaction_id` varchar(256) NOT NULL,
    `amount` decimal(10,2) DEFAULT NULL,
    PRIMARY KEY (`transaction_id`),
    KEY `source_account_id` (`source_account_id`),
    KEY `destination_account_id` (`destination_account_id`),
    CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`source_account_id`) REFERENCES `accounts` (`id`),
    CONSTRAINT `transactions_ibfk_2` FOREIGN KEY (`destination_account_id`) REFERENCES `accounts` (`id`)
);
```

### Navigate to project folder and create .env file

```
touch .env
echo DB_HOST=localhost >> .env
echo DB_PORT=3306 >> .env
echo DB_USERNAME=<DB_USERNAME> >> .env
echo DB_PASSWORD=<DB_PASSWORD> >> .env
echo DB_DATABASE=<DB_DATABASE> >> .env
```

### Start the application
```
go run .
```

## APIs
List of available endpoints.

### Postman
Import the "Internal Transfers.postman_collection.json" file into Postman to get the list of available endpoints.

### Accounts
`GET /accounts/{id}`

Path params:
- `id`

Response Body:
```
{
    "account_id": 1,
    "balance": "96.5"
}
```

`POST /accounts`

This endpoints takes two fields, `account_id` and `initial_balance`. `account_id` should be a `int` value, whereas `initial_balance` is a `string` value wrapping the amount value. Allows up to 2 decimal places, any decimal values after that are truncated.

Request Body:
```
{
    "account_id": 11,
    "initial_balance": "100.20"
}
```

Response Body:
```
Account created
```

### Transactions
`POST /transactions`

This endpoint allows two accounts to transfer money from each other. `source_account_id` and `destination_account_id` should both be present in the `accounts` table.

`transaction_id` is supposed to be a unique identifier for the transaction. This is meant to prevent duplicate transactions from occurring, if say there are multiple upstream requests to make the transaction. The example below uses a timestamp format, but it should be something more unique. One way would be to concatenate the timestamp with some metadata of the two accounts involved to generate a more unique identifier. Another way is to generate a random string of characters to be used as an identifier. Trying to insert the same `transaction_id` will fail since it's used as a primary key in the `transactions` table, rolling back any changes made by the transaction so far.

Request Body:
```
{
    "source_account_id": 2,
    "destination_account_id": 1,
    "transaction_id": "20240425151400",
    "amount": "1.50"
}
```

Response Body:
```
Transaction completed successfully
```

## Testing
Integration testing should be done for handling different scenarios but due to lack of time, this is a TO-DO.

## TO-DO
Would want to make the exception handling more elegant and extensible.

The setup for the DB tables could also be handled better using a DB migration tool/library, so that addition of new tables and setting up of the local environment would be easier. 

More tests to be written as well, both unit and integration tests.

To deploy this application on Docker would be a nice-to-have as well.
