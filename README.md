# Pismo

## Developing locally
#### Setting up the database
Run `sh ./run` to start the service and mysql db. Hard coded sample data will be generated. If you want to clear any data you have added and restart to the sample data, run `sh ./run-clean`.  
If you are on a Windows machine, you can run the commands in the script manually.

#### MySQL access
Run `mysql -h 127.0.0.1 -P 3306 -u root -p` to log into the mysql cli. User name and password are `root` for testing purposes.

#### IDE
VS Code was used to develop this app, so the `launch.json` is already configured. If you are using an alternate ID, you will need to set up your own build configuration.

## API Contract
### Base URL
```bash
http://localhost:8080
```

Endpoints
1. Get an Account by ID
- URL: `/accounts/{id}`
- Method: GET
- Description: Retrieves account details by the given `id`.
- Path Parameters:
    - `id` (integer) - The ID of the account to retrieve.
- Response:

    - Status Code: 200 OK

        ```json
        {
            "id": 1,
            "document_number": "123456789"
        }
        ```
    - Status Code: 404 Not Found

        ```json
        {
            "error": "Account not found"
        }
        ```
    - Status Code: 400 Bad Request

        ```json
        {
            "error": "Invalid account ID: <account_ID>"
        }
        ```
    - Status Code: 500 Internal Server Error

        ```json
        {
            "error": "Internal server error"
        }
        ```

2. Create an Account
- URL: `/accounts`
- Method: POST
- Description: Creates a new account with the provided document number.
- Request Body:

    ```json
    {
        "document_number": "string" 
    }
    ```
- Response:
    - Status Code: 201 Created
    - Body:

        ```json
        {
            "message": "Successfully created new account with ID <account_ID>"
        }
        ```
    - Status Code: 400 Bad Request

        ```json
        {
            "error": "No document number provided"
        }
        ```
    - Status Code: 500 Internal Server Error

        ```json
        {
            "error": "an account with that document number already exists"
        }
        ```

3. Create a Transaction
- URL: /transactions
- Method: POST
- Description: Creates a new transaction for the specified account.
- Request Body:

    ```json
    {
        "account_id": 1,
        "operation_type_id": 4,
        "amount": 100.50,
        "event_date": "2024-09-17T15:04:05Z" // optional, will default to current timestamp
    }
    ```
- Response:
    - Status Code: 201 Created

        ```json
        {
            "message": "Successfully created new transaction with ID <transaction_ID>"
        }
        ```
    - Status Code: 400 Bad Request

        ```json
        {
            "error": "Invalid request payload"
        }
        ```
    - Status Code: 500 Internal Server Error

        ```json
        {
            "error": "Internal server error"
        }
        ```

## Notes
- `operation_type_id`: Represents the type of operation:  
    - `1`: Normal Purchase (Debit)  
    - `2`: Purchase with Installments (Debit)  
    - `3`: Withdrawal (Debit)  
    - `4`: Credit Voucher (Credit)  
- `amount` should be positive for credits and negative for debits.

## Auth
- TODO...

## Content Type
- All request and response bodies must use `Content-Type: application/json`.

## Example cURL Requests
#### Create Account
```bash
curl -X POST "http://localhost:8080/accounts" \
  -H "Content-Type: application/json" \
  -d '{"document_number":"123456789"}'
```
#### Get Account by ID
```bash
curl -X GET "http://localhost:8080/accounts/1" \
  -H "Content-Type: application/json"
```
#### Create Transaction
```bash
curl -X POST "http://localhost:8080/transactions" \
  -H "Content-Type: application/json" \
  -d '{"account_id":1,"operation_type_id":4,"amount":100.50}'
```

## Sample Tables
```
Accounts
+------------+-----------------+
| account_id | document_number |
+------------+-----------------+
|          1 | 12345678900     |
|          2 | 12345678901     |
|          3 | 12345678902     |
|          4 | 12345678903     |
+------------+-----------------+

OperationTypes
+-------------------+----------------------------+
| operation_type_id | description0               |
+-------------------+----------------------------+
|                 1 | Normal Purchase            |
|                 2 | Purchase with installments |
|                 3 | Withdrawal                 |
|                 4 | Credit Voucher             |
+-------------------+----------------------------+

Transactions
+----------------+------------+-------------------+--------+---------------------+
| transaction_id | account_id | operation_type_id | amount | event_date          |
+----------------+------------+-------------------+--------+---------------------+
|              1 |          1 |                 1 | -50.00 | 2020-01-01 10:32:08 |
|              2 |          1 |                 1 | -23.50 | 2020-01-01 10:48:12 |
|              3 |          1 |                 1 | -18.70 | 2020-01-02 19:01:23 |
|              4 |          1 |                 4 |  60.00 | 2020-01-05 09:34:19 |
+----------------+------------+-------------------+--------+---------------------+
```