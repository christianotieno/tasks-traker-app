# Maintenance Task Tracker

## Description
This application is a software solution designed to account for maintenance tasks performed during a working day. It provides functionality for two types of users: Managers and Technicians. The Technicians can create, update, and view their own tasks, while the Managers have access to all tasks, including the ability to delete them. Additionally, Managers receive notifications whenever a Technician performs a task.

## Features Implemented
The following features have been implemented:
1. Create API endpoint to save a new task: The application provides an API endpoint that allows Technicians to save a new task. The endpoint accepts the task details, such as summary and date, and stores them in the MySQL database.
2. Create API endpoint to list tasks: The application also includes an API endpoint to retrieve a list of tasks. Technicians can view their own tasks, while Managers can view all tasks.
3. Local development environment using Docker: The application is containerized using Docker, which provides a consistent and isolated development environment. The Docker setup includes the service itself and a MySQL database, ensuring that the application can be easily set up and run on different systems.
4. MySQL database for data persistence: The application utilizes a MySQL database to store task data. This ensures that task information is securely stored and can be accessed efficiently when needed. 
5. Unit tests: All implemented features have accompanying unit tests to ensure their proper functioning. These tests verify the behavior of the endpoints, data storage, and notification logic.
6. Tilt for local development: The application utilizes Tilt to facilitate local development. Tilt provides a convenient way to build, run, and test the application. It also provides a way to view the application logs and access the MySQL database.
7. Kafka for notification: The application utilizes Kafka to facilitate notification. Kafka provides a reliable and scalable way to send notifications to the Manager. It also ensures that the notification process is decoupled from the application flow, allowing the application to continue functioning smoothly.
8. Kubernetes deployment: The application includes Kubernetes object files to facilitate deployment. These files can be used to deploy the application, the MySQL database, and any required dependencies. This ensures that the application can be easily scaled, managed, and maintained in a production environment.

## Futures currently in Development & Testing

In the coming days, the following tasks are planned to be accomplished:

1. Message broker integration: Currently, the notification logic is implemented as a simple print statement. To enhance scalability and decouple the notification process from the application flow, a message broker will be integrated. This will enable asynchronous communication between components, ensuring that the Manager receives notifications reliably and without blocking the application.
2. Kubernetes deployment: Kubernetes object files will be created to facilitate the deployment of this application. By utilizing Kubernetes, the application can be easily scaled, managed, and maintained in a production environment. The Kubernetes deployment files will be provided, including configuration for deploying the service, the MySQL database, and any required dependencies.

## Future Improvements

- Integrate more tests and improve test coverage: Currently, the application has unit tests for the implemented features. In the future, more tests will be added to ensure that the application is functioning as expected. Additionally, the test coverage will be improved to ensure that all code paths are tested.
- Swagger documentation: Currently, the application does not have any documentation. In the future, Swagger documentation will be implemented to provide a detailed description of the API endpoints. This will allow for easier debugging and troubleshooting of the application.
- Implement industry-standard authentication and authorization: Currently, the application only have a very basic authentication or authorization. In the future, these features will be implemented to ensure that only proper authorized users can access the application.
- Implement a frontend: Currently, the application only has a backend. In the future, a frontend will be implemented to provide a user interface for the application. This will allow users to interact with the application without having to use the API endpoints directly.

## Future Improvements
The following features are planned to be implemented in the future:

- Implement a CI/CD pipeline: Currently, the application does not have a CI/CD pipeline. In the future, a pipeline will be implemented to ensure that the application is built, tested, and deployed automatically. This will ensure that the application is always in a working state and that new features can be deployed quickly and reliably.
- Implement a logging solution: Currently, the application does not have a logging solution. In the future, a logging solution will be implemented to ensure that the application logs are stored and can be accessed when needed. This will allow for easier debugging and troubleshooting of the application.
- Implement a monitoring solution: Currently, the application does not have a monitoring solution. In the future, a monitoring solution will be implemented to ensure that the application is always running and performing as expected. This will allow for easier debugging and troubleshooting of the application.
- Implement a metrics solution: Currently, the application does not have a metrics solution. In the future, a metrics solution will be implemented to ensure that the application performance can be monitored and analyzed. This will allow for easier debugging and troubleshooting of the application.
- Implement a security solution: Currently, the application does not have a security solution. In the future, a security solution will be implemented to ensure that the application is secure and protected from malicious attacks. This will allow for easier debugging and troubleshooting of the application.
- Implement a data backup solution: Currently, the application does not have a data backup solution. In the future, a data backup solution will be implemented to ensure that the application data is backed up and can be restored when needed. This will allow for easier debugging and troubleshooting of the application.
- Implement a data recovery solution: Currently, the application does not have a data recovery solution. In the future, a data recovery solution will be implemented to ensure that the application data can be recovered when needed. This will allow for easier debugging and troubleshooting of the application.
- Implement a data migration solution: Currently, the application does not have a data migration solution. In the future, a data migration solution will be implemented to ensure that the application data can be migrated when needed. This will allow for easier debugging and troubleshooting of the application.

## Technologies Used
- Go
- MySQL
- Docker
- Tilt
- Kafka

## Setup and Installation
To set up the application locally, follow these steps:

- Clone the repository: `git clone git@github.com:christianotieno/tasks-traker-app.git`
- Navigate to the project directory: cd maintenance-task-tracker 
- Run the app using `go run main.go`
- The application should now be accessible at http://localhost:8000

## Usage
To use the application, follow these steps:

- Create a new manager by sending a POST request to http://localhost:8000/users with the following payload:
```
{
    "first_name": "Manager first name",
    "last_name": "Manager last name",
    "email": "Your email",
    "password": "Your password",
```

- Create a new technician by sending a POST request to http://localhost:8000/users with the following payload:
```
{
    "first_name": "Technician first name",
    "last_name": "Technician last name",
    "email": "Your email",
    "password": "Your password",
    "manager_id" "manager id you want associated with the technician
```

You will get a token in the response. Copy the token and use it in the next step.

- Login by sending a POST request to http://localhost:8000/login with the following payload:
```
{
    "email": "Your email",
    "password": "Your password",
}
```

You will get a token in the response. Copy the token and use it in the next step.


- Paste the token in the Authorization header as a Bearer token. You can now access the protected endpoints.
- If you are a technician, you can create a task by sending a POST request to http://localhost:8000/tasks with the following payload:
```
{
    "summary": "Task summary",
    "date": "2021-10-10T10:10:10Z",
    "technician_id": "Technician id you want to assign the task to"
}
```
- Your privileges will be checked to ensure that you are allowed to create a task and update the task the specified technician. If you are not allowed, you will get an error message. You cannot delete a task, only your manager can do that.

- As a manager, you can view tasks of your technicians by sending a GET request to http://localhost:8000/tasks, there are plans to move this endpoint to  http://localhost:8000/tasks?manager_id=1
- You can also delete a task by sending a DELETE request to http://localhost:8000/tasks/1
- You can only delete taks that are assigned to your technicians. You cannot delete a task that is assigned to another manager's technician.

These are the endpoints that are currently available:

- Create a new task by sending a POST request to http://localhost:8000/tasks.
- Update a task by sending a PUT request to http://localhost:8000/tasks/{id}
- List all tasks by sending a GET request to http://localhost:8000/tasks
- List all tasks for a specific technician by sending a GET request to http://localhost:8000/users/{id}/tasks
- Delete a task by sending a DELETE request to http://localhost:8000/tasks/1
- View all tasks completed by a technician by sending a GET request to http://localhost:8000/tasks/technicians/1/completed
- View all tasks by all technicians for a specific manager by sending a GET request to http://localhost:8000/tasks

### Contributing

If you encounter any issues or have suggestions for enhancements, please submit an issue or a pull request on the repository.