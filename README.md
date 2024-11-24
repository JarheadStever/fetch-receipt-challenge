# fetch-receipt-challenge
My solution / implementation for Fetch's receipt-processor-challenge

## TL;DR

This webservice fulfills the requirements set out by Fetch in this [set of instructions](./instructions/prompt.md) and matches [this provided OpenAPI spec](./instructions/apispec.yml)


#### To run:
- Clone this repository
```
git clone https://github.com/JarheadStever/fetch-receipt-challenge
```
- Navigate to the project directory and build the app
```
cd fetch-receipt-challenge
go build -o jaredsapp
```
- Run it
```
./jaredsapp
```
NOTE: The default port is `:3005`. To specify a different port, use the `--port` flag like so:
```
./jaredsapp --port 8080
```

To try it out, `POST` to `localhost:3005/receipts/process` this JSON body:
```
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```

You should receive a UUID for that receipt and can retrieve its score by `GET`ting `localhost:3005/receipts/{id}/points`. If everything is working, you should see a score of `28`.


## Thoughts
***A note on Validation:*** I think that ultimately, the way I abstracted the validation a bit into a more applicative instead of monadic style is a good approach and would be beneficial for testing/logging. That said, it also made things a bit more complex than I think was necessary for this exercise. The simple validation (see the original code here) also met the requirements and I was torn on which route to go for my submission. Interesting to think about...

####  Things I like about this implementation
- It's simple. No overengineering here, just a few simple files that perform the expected operations to calculate, store, and fetch receipt scores.
- It's lightweight. Our dependencies are minimal and we don't rely on many libraries to handle things for us. There are libraries I could've used for things like checking for alphanumeric characters or converting strings to various other types, but for the most part I tried to make things as readable and easy to understand as possible. Yay for code stewardship! (And yay for avoiding as many of the inevitable CVEs that come over time with those libraries as possible!)
- It's safe. We have validation in place for receipts and nothing is stored other than a UUID and a corresponding int. If we had a true SQL database or something of that nature, we'd need to take extra precautions to prevent things like SQL injection and other vulnerabilities. This solution won't scale, but it simplifies the security aspect of the API.

####  Things I don't like about this implementation
- Our file/project structure isn't exactly "proper". This builds, runs, and works as intended, but, for example, we don't have the traditional "pkg", "cmd", etc. layout in this project. There is also probably a more "correct" way to organize the contents/naming of these files, but important conventions were followed, and everything makes enough logical sense. For the sake of meeting requirements here as simply and efficiently as possible, there were conscious decisions I made, but if this code were to ever scale or be part of a larger codebase, there'd need to be some more structure.
- There isn't a ton of consistency with exported vs un-exported naming here. I left most things lowercase since this is a single package/project and it's not getting used anywhere else.
- More in-depth commenting/docstrings never hurts (within reason, of course).

####  Additional considerations/ideas
- Testing is always good, but writing unit tests seemed to be out of scope for this project. I tested things locally until I was satisfied, as will whoever evaluates this, I'm sure. So, just acknowledging the lack of tests here.
- The barebones `map[UUID]int` storage implementation works, but it certainly wouldn't scale. For the sake of this exercise, it's sufficient, but it would be neat to implement SQLite or something of that nature to have a lightweight db and a more formal data storage solution.
- A Dockerfile to build/run this service would be cool too, but again seemed out of scope. Instructions said that Go was the preferred language and "If you are not using Go, include a Dockerized setup to run the code" which led me to believe a Dockerized solution was not necessary if using Go.
- More verbose logging / http request logging would be a cool addition as well with another flag. Would be helpful for further development/scaling, but again was out of scope.
