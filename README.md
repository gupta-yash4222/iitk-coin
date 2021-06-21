# iitk-coin

## Dependencies of the Project 
1. "github.com/mattn/go-sqlite3"
2. "github.com/gorilla/mux"
3. "golang.org/x/crypto/bcrypt"
4. "github.com/dgrijalva/jwt-go"

To get all these dependencies type the command `go get <package-name>` in terminal or powershell, or simply type `go mod tidy` in terminal or powershell

## Instructions on Sending requests to various end-points

### `/signup` endpoint (POST)
You are required to send a request in JSON format as - 

    {
        "rollno": <roll-number>
        "name": <name>
        "password": <password>
    }

The password must contain at least one character, otherwise an error will come. 

### `/login` endpoint (POST)
You are required to send a request in JSON format as -

    {
        "rollno": <roll-number>
        "name": <name>
        "password": <password>
    }

Secret Key used for generating JWT : **CROWmium** (this can be edited in `.\server\login.go`)  

Duration after which JWT expires : **25sec** (this can be edited in `.\server\login.go` by changing the value of `validDuration` variable)

### `/secretpage` endpoint (GET)
In my implementation, I am embedding the obtained JWT in a cookie and thus in every subsequent request the browser will send that cookie. Thus, this endpoint only requires you to send a blank request.
While using *Postman*, enable the *Interceptor* to entertain cookies. 

The endpoint can only be accessed by a user who has admin rights. In my case, I have given these rights to myself with credentials - 

    {
        "rollno": 190998
        "name": "Yash Gupta"
        "password": "hsay"
    }

For any other user, *Access Denied* will come. 
Currently, I have hard-coded this admin-rights thing, but later I will include that in my database for more smooth functioning.  
 

### `/getUsers` endpoint (GET)
This endpoint returns all the registered users and their details, but not their password (although even if it had been returning them, they would have been of no use :) ).
To entertain this endpoint too, user just requires to send a blank request. The output(response) will look like - 

    Rollno: <roll-number>, Name: <name>
    ...

### `/rewardCoins` endpoint (POST)
This endpoint is to add given number of coins in the specified users's account. This endpoint is primarily aimed to be only used by the admin of `iitk-coin`, but in the present implementation this endpoint can be entertained by anyone. 
For this endpoint, you require to send a JSON request as - 

    {
        "rollno": <roll-number>
        "coins": <no-of-coins>
    }

If the rollno is registered in the database then a result saying **Transaction successful** will come. 

### `/getBalance` endpoint (GET)
This endpoint returns the current balance (no of coins) of the specified user's account. You need to send a JSON request as -

    {
        "rollno": <roll-no>
    }

If the rollno is registered in the database, then you will get a message showing your current balance. 

### `/transferCoins` endpoint (POST)
This endpoint is for transfering coins between 2 users. For sending some number of coins from one user to another user, you will need to use this endpoint. You need to send a JSON request as -

    {
        "senderRollno": <sender-rollno>
        "receiverRollno": <receiver-rollno>
        "coins": <no-of-coins>
    }

First the backend checks if both the sender and receiver are registered in the database, and if not then it will return an error message. If sender has sufficient balance, then the backend deducts the given number of coins from the sender's account (if it fails then the whole transaction is aborted) and then add that number of coins into the receiver's account (if it fails then the whole transaction is aborted). Also, while doing a single transaction, the backend acquires a lock on the whole database and thus any other transaction cannot be entertained concurrently. 

### `/hihi` endpoint
This is just for fun :grin: 
