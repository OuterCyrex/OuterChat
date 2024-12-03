# OuterChat

## Overview

OuterChat is a Go-based application developed using the Gin framework, designed to provide robust API endpoints for user management and real-time messaging capabilities. The application leverages WebSocket for real-time communication and Redis for publish-subscribe messaging patterns. Currently, the project is under active development, with additional modules for group and message management planned.

## Key Features

- **User Management**: RESTful API endpoints for managing user data.
- **Real-Time Messaging**: WebSocket integration for instant messaging features.
- **Publish-Subscribe Model**: Utilizes Redis for efficient message distribution.

## API Documentation

The API documentation is accessible via the Swagger UI, which can be viewed at the `/swagger` endpoint.

### Redis Integration

The application uses Redis for publish-subscribe functionality. This allows for efficient and scalable messaging across the system.

### WebSocket Integration

WebSocket is used for real-time communication, enabling features like instant messaging within the application.

## Roadmap

- Implement group management module.
- Develop message handling module.
- Enhance security features.
- Improve error handling and logging.

## What I want to say

This is just a short Instant Message Project for Studying Message Queue as well as Websocket so there's gonna be lots of mistake, don't be too serious :)

## License

[MIT](https://choosealicense.com/licenses/mit/)
