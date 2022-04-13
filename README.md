# Simple URL Shortener in Go

## Description
The project supports:
- Creating shortened urls for valid input url strings. 
- Fetching content based hosted at the original url using the shortened url string.

## Running the application 
  - docker-compose up & 
  - go run main.go
  - By default applications will be available at localhost:8185

## API Usage
- Create shortened URL : 
  - resource path: /url-shortener
  - method: POST 
  - request payload
    ```json
      {
        "originalUrl": "<original url to be shortened>",
        "requesterId": "<an identifier for requester>",
        "ExpiryDate": "<when the shortened url should exoiry>"
      } 
      ```       
  - response codes: [201,400]
  - response payload:
     ```json
      {
        "id": "unique id of the mapping between original url and shortened url",
        "originalUrl": "<the original url>",
        "expiryDate": "<expiry date for the mapping>",
        "shortenedUrl": "<the actual shortened url>",
        "createdDate": "<insert date of the mapping>",
        "requesterId": "<an identifier for requester>",
        "content": "<actual site content histed at originalUrl>"
     }
     ```   
   - example call:  
     ```             
     curl --location --request POST 'http://localhost:8185/url-shortener' \
        --header 'Content-Type: application/json' \
        --data-raw '{
            "originalUrl" :"https://jsonplaceholder.typicode.com/",
            "requesterId": "sunny",
            "ExpiryDate": "2023-04-12T02:05:55Z"
        }'
       ```
                     
             
    - example response:
      ```json
      {
          "id": "4",
          "originalUrl": "https://jsonplaceholder.typicode.com/",
          "expiryDate": "2023-04-12T02:05:55Z",
          "shortenedUrl": "TmeFhogLWy",
          "createdDate": "2022-04-13T12:48:11Z",
          "requesterId": "sunny",
          "content": "\n<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta lesheet\" href=\"/style.css\....and so on"       }
     ```
 - Get shortened url :
    - resource path: /url-shortener/{shortUrlString}
    - method: GET   
    - response codes: [200,404]
    - response payload:
     ```json
      {
        "id": "unique id of the mapping between original url and shortened url",
        "originalUrl": "<the original url>",
        "expiryDate": "<expiry date for the mapping>",
        "shortenedUrl": "<the actual shortened url>",
        "createdDate": "<insert date of the mapping>",
        "requesterId": "<an identifier for requester>",
        "content": "<actual site content histed at originalUrl>"
     }
     ```   
    - example call:  
     ```             
     curl --location --request GET 'http://localhost:8185/url-shortener/TmZcbm7O78'
      ```
                     
             
    - example response:
      ```json
      {
       "id": "1",
       "originalUrl": "https://jsonplaceholder.typicode.com/",
       "expiryDate": "2023-04-12T02:05:55Z",
       "shortenedUrl": "TmZcbm7O78",
       "createdDate": "2022-04-13T10:06:11Z",
       "requesterId": "sunny",
       "content": "\n<!DOCTYPE html>\n<html lang=\"en\">\n<head....and so on"
       }
     ```
## Testing
   - Curently only integrations tests ara available in the integration_tests folder
   - Once application is running change directory to <project-root-folder>/integration_tests and execute:
      ```
         newman run URL_shortener_tests.postman_collection.json
      ```
  
 ## Possible future enhancements
    - Add api, service and store tests.
    - Consider fetching shortened url from dedicated unique service.
    - Improve performance, current response time ~150-600ms. 
