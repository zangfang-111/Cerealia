# Browser frontend

## Setting up

Requirements:

1. yarn >= 1.12
1. nodejs >= 8.0
1. flow-typed >= 2.5

Installation:

1. `yarn install --dev --prefer-offline`  (or `make install-deps`)

Configuration: copy end edit the .env.example file

    cp .env.example .env
    $EDITOR .env


## Running

To start development server:

1. `cp .env.example .env`
1. modify .env and set `REACT_APP_API_HOST` to your server location
1. `yarn start`

To build for production:

1. `cp .env.example .env`
1. modify .env and set `REACT_APP_API_HOST` to your server location
1. `yarn run build`


### TESTING

    yarn run test

## Checking out the code

Before checking out the code make sure your code follows all Cerealia standard, among others:

+ passes linting: `yarn run lint`
+ passes all tests: `test:all`
