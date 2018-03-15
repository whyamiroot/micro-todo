# Micro TODO admin web application

Micro-TODO-admin shows all services, registered in the service registry. It keeps the list up to date by pinging health 
route of each service individually and refreshing services list from registry.

*TODO*: 
- collect logs from `Logger` service
- implement some basic authentication

## How to build

1. Go to `micro-todo-admin` directory
2. Install dependencies with `npm install`
3. Run dev application with `npm run dev` or release application with `npm run build`