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

### `/hihi` endpoint
This is just for fun :grin: 
