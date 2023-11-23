# Hotel Reservation Backend

## Overview

This backend system is designed to manage hotel reservations through a clean and organized architecture. The system follows a structured flow from handling routes to interacting with datastores, promoting maintainability and scalability.

## Architecture

- **Routing and Resource Handling:**
  - Each route is handled by a dedicated resource handler.
  - Resource handlers are responsible for invoking specific methods on the manager in the business layer.

- **Business Layer:**
  - The business layer, managed by the manager, acts as an intermediary between resource handlers and datastores.
  - Managers encapsulate business logic and provide a clear API that mirrors the business use cases.

- **Datastores:**
  - Datastores store and retrieve data, and they are accessed exclusively through the manager in the business layer.
  - Resource handlers should avoid direct calls to the database to maintain a separation of concerns.

## Best Practices

- **Consistent API:**
  - Ensure that the manager's API closely resembles the business use cases to enhance clarity and maintainability.

- **Usecase-Centric Approach:**
  - Implement each business use case as a separate method on the manager. This promotes a modular and organized codebase.

- **Avoid Direct Database Calls:**
  - Resource handlers should refrain from making direct calls to the database. All database interactions should be routed through the manager.


## Guidlines 
### Logging: REF [Design Philosophy On Logging]((https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html))

- Logging should be isolated to the single purpose of debugging errors.
- Logging and error handling are coupled and not separate concerns. 
- Handling an error means:
  - The error has been logged.
  - The application is back to 100% integrity.
  - The current error is not reported any longer.
- Packages that are reusable across many projects only return root error values.
- If the error is not going to be handled, wrap and return up the call stack.
- Once an error is handled, it is not allowed to be passed up the call stack any longer.



## Getting Started

Please look at the `Taskfile.yaml` in the root directoty. The project is entirlt managed using task.

## Dependencies

- mongodb

## Dev Environment

- Visual Studio Code
- Postman
- MongoDB Compass

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute the code in accordance with the terms of the license.

