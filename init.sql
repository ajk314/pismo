CREATE TABLE IF NOT EXISTS Accounts (
    account_id INT AUTO_INCREMENT PRIMARY KEY,
    document_number VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS OperationTypes (
    operation_type_id INT AUTO_INCREMENT PRIMARY KEY,
    description0 VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS Transactions (
    transaction_id INT AUTO_INCREMENT PRIMARY KEY,
    account_id INT,
    operation_type_id INT,
    amount DECIMAL(10, 2),
    event_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES Accounts(account_id),
    FOREIGN KEY (operation_type_id) REFERENCES OperationTypes(operation_type_id)
);

INSERT INTO Accounts (account_id, document_number) 
VALUES (1, '12345678900');

INSERT INTO OperationTypes (operation_type_id, description0)
VALUES 
(1, 'Normal Purchase'),
(2, 'Purchase with installments'),
(3, 'Withdrawal'),
(4, 'Credit Voucher');

INSERT INTO Transactions (transaction_id, account_id, operation_type_id, amount, event_date)
VALUES 
(1, 1, 1, -50.0, '2020-01-01 10:32:07.719922'),
(2, 1, 1, -23.5, '2020-01-01 10:48:12.2135875'),
(3, 1, 1, -18.7, '2020-01-02 19:01:23.1458543'),
(4, 1, 4, 60.0, '2020-01-05 09:34:18.5893223');