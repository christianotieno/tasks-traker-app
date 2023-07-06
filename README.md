# Maintenance Task Tracker
## Description


This application is a software solution designed to account for maintenance tasks performed during a working day. It provides functionality for two types of users: Managers and Technicians. The Technicians can create, update, and view their own tasks, while the Managers have access to all tasks, including the ability to delete them. Additionally, Managers receive notifications whenever a Technician performs a task.

## Features Implemented

1. Create API endpoint to save a new task: The application provides an API endpoint that allows Technicians to save a new task. The endpoint accepts the task details, such as summary and date, and stores them in the MySQL database.
2. Create API endpoint to list tasks: The application also includes an API endpoint to retrieve a list of tasks. Technicians can view their own tasks, while Managers can view all tasks.
3. Notification of task performed by Technician: Whenever a Technician performs a task, the Manager is notified with a message. This notification is currently implemented as a print statement. The notification does not block any HTTP requests and allows the application to continue functioning smoothly.
4. Local development environment using Docker: The application is containerized using Docker, which provides a consistent and isolated development environment. The Docker setup includes the service itself and a MySQL database, ensuring that the application can be easily set up and run on different systems.
5. MySQL database for data persistence: The application utilizes a MySQL database to store task data. This ensures that task information is securely stored and can be accessed efficiently when needed. 
6. Unit tests: All implemented features have accompanying unit tests to ensure their proper functioning. These tests verify the behavior of the endpoints, data storage, and notification logic.


## Future Development
In the coming days, the following tasks are planned to be accomplished:

1. Message broker integration: Currently, the notification logic is implemented as a simple print statement. To enhance scalability and decouple the notification process from the application flow, a message broker will be integrated. This will enable asynchronous communication between components, ensuring that the Manager receives notifications reliably and without blocking the application.

2. Kubernetes deployment: Kubernetes object files will be created to facilitate the deployment of this application. By utilizing Kubernetes, the application can be easily scaled, managed, and maintained in a production environment. The Kubernetes deployment files will be provided, including configuration for deploying the service, the MySQL database, and any required dependencies.


## Setup and Installation
To set up the application locally, follow these steps:

- Clone the repository: `git clone git@github.com:christianotieno/tasks-traker-app.git`
- Navigate to the project directory: cd maintenance-task-tracker 
- Run the app using `go run main.go`
- The application should now be accessible at http://localhost:8000


### Contributing

If you encounter any issues or have suggestions for enhancements, please submit an issue or a pull request on the repository.

### License
This project is licensed under the MIT License.