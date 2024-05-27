# Tabungan API Go

## Run APP
1. Clone this repository
2. Copy file .env.example to .env
3. Run using command below
```
docker-compose up
```

## Run Test
```
go test -v 
```

## Docs 
For docs access on this link localhost:port/api/v1/docs/index.html

## Following endpoints:

| Method | Route     | Header                               |  Body                                                       |
| ------ | --------- | ------------------------------------ | ----------------------------------------------------------- |
| POST   | /daftar   |                                      |  `{"name": "Dian", "phone_number": "0821", "pin": "1234" }` |
| POST   | /tabung   | Basic Auth, using account_no dan pin |  `{"amount": 50000}`                                        |
| POST   | /tarik    | Basic Auth, using account_no dan pin |  `{"amount": 50000}`                                        |

* Note: 
* Basic AUth can generated online using this link https://www.debugbear.com/basic-auth-header-generator,
* fil username using account_no, and password using pin.
* result generated copy to header. Ecample -> Authorization: Basic am9obkBleGFtcGxlLmNvbTphYmMxMjM=
