# URL Shortener

## 1. Create .env file        
        
How .env file looks:        
        
#Server settings        
SERVER_PORT=1323        
        
#Redis settings        
REDIS_ADDRESS=                        
REDIS_PASSWORD=                         
HASH_NAME=                      # hash where you want to store the data.         

## 2. HTTP Requests


### POST /cutURL

json
{
    "longURL": "https://www.example.com/very/long/url",        
    "expirationTime": 12                             // Not used now
}

Returns short URL, for example "s4d5sd"

### GET "/SHORT_URL"

If SHORT_URL exists, it will redirect you to the corresponding long URL
