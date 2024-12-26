# URL Shortener

This is a simple URL shortener application built with Go. It allows users to shorten long URLs and redirect to the original URLs using the shortened codes.
The web-app provides authentication mechanism using which user can login to the web application. After the user is authenticated, user will be able to call protected API endpoint which generates short-URL.
Calling the short-URL is public, and it's not protected. 

Short codes are randomly generated for long-URL using base-32 character set of length 7-characters.

## Features

- Shorten long URLs
- Redirect to original URLs using shortened codes

## Prerequisites

- Go 1.19 or later
- ZITADEL account for authentication/authorization
- Setup authentication for the web-app in Zitadel console, similar to what described [here](https://zitadel.com/docs/examples/login/go)
- Also setup authorization for the APIs in Zitadel console, similar to what described [here](https://zitadel.com/docs/examples/secure-api/go)

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/sankalp-r/url-shortner.git
    cd url-shortner
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

## Running the Application

1. Build and run the application:

   Set the following variables in the Makefile:
    - DOMAIN, KEY, CLIENT_ID, REDIRECT_URI (refer [this](https://zitadel.com/docs/examples/login/go#set-up-application)))
    - KEY_FILE (refer [this](https://zitadel.com/docs/examples/secure-api/go#set-up-application-and-obtain-keys))

    
Then execute:

 ```sh
 make run
 ```

2. The application will be available at `http://localhost:8080`.

## API Endpoints

- `POST /v1/short`: Shorten a long URL. This is a protected endpoint. You must fetch the access token from Zitadel account before calling this API.
  ```sh
  curl --location 'http://localhost:8080/v1/short' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <api-token>' \
  --data '{"url":"https://www.google.com"}'
  ```
- `GET /v1/{shortURLCode}`: Redirect to the original URL using the shortened code

## Usage

1. Open `http://localhost:8080` in your browser.
2. After performing login, enter a long-URL in the text box and click the "Shorten" button.
3. The shortened URL will be displayed.
4. Copy the shortened URL and open it your browser, it will redirect to original long URL.


## Enhancement scope
1. This is a very simple web-application, and it's not production ready. User-interface can be improved.
2. In this implementation, short-URL code generation and storage are handled in-memory which is not durable.
   Database can be introduced for providing durability by extending `storage` interface.
3. Rate-limiter can be introduced for prevention from attackers.

## References

- [Gorilla Mux](https://github.com/gorilla/mux)
- [ZITADEL](https://zitadel.com)
- [ZITADEL Go Web-app](https://zitadel.com/docs/examples/login/go)
- [ZITADEL Go API](https://zitadel.com/docs/examples/secure-api/go)
