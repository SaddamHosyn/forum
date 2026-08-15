# Forum-Go Project

Reel Movie Talk is an online forum dedicated to movie discussions, where users can share their opinions and engage in conversations about films. <br>
It's a platform to connect, debate, and explore the world of cinema together. <br><br>
The Reel Movie Talk project focused on registration, login, session management, database interactions (SQLite), Docker setup, and functional aspects of the forum. <br>
Users can register and login to create posts, and comments. <br>
Registered users can also interact with content by liking and disliking posts, and browse posts by genre or movie.


## Table of Contents
- [Description](#description)
- [Features](#Features)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Requirements](#requirements)
- [Learning Objectives](#learning-objectives)
- [Last Updated](#last-updated) 
- [Authors](#authors) <br><br>


[Back To The Top](#forum-go-project) 

## Description
<br>
This project entails the development of a web-based forum application that facilitates user interaction and content sharing. <br><br>
The forum is designed to support user registration, authentication, and communication through posts and comments. <br>
It incorporates features such as category association, content filtering, and a like/dislike system. <br><br>

The application is built using Go for the backend, with SQLite as the database management system. <br>
It employs standard web technologies including HTML, HTTP, and cookies for session management. <br>
The entire application is containerized using Docker to ensure consistent deployment across different environments. <br><br>
This forum project serves as a comprehensive exercise in full-stack web development,
covering aspects from database management to user interface design, <br>
while also addressing important considerations such as security and scalability.

<br>

[Back To The Top](#forum-go-project) 

## Features

**Key aspects of the project include:** <br>

- Implementation of user authentication and session management
- Database design and query optimization for efficient data handling
- Development of a responsive and intuitive user interface
- Integration of content filtering and categorization features
- Application of security best practices, including password encryption

<br>

[Back To The Top](#forum-go-project) 

## Project Structure

```text
forum/
├── main.go               # Entry point (initializes DB & starts server)
├── assets/               # Static assets (CSS, JS, Images)
├── auth/                 # Password hashing & user auth helpers
├── database/             # SQLite connection, schema & query functions
├── handler/              # HTTP endpoint route handlers
├── middleware/           # Session management & CORS middlewares
├── model/                # Data structures (User, Post, Comment, Category)
├── pkg/utils/            # Input validation & utility functions
├── render/               # Template parsing engine (render.go)
├── server/               # Router registration & HTTP server setup
├── templates/            # HTML view templates
├── tests/                # Consolidated unit test suite (auth, validation, handlers)
├── Dockerfile            # Multi-stage Docker build config
├── docker-compose.yml    # Container orchestration setup
└── .dockerignore         # Docker context exclusion rules
```

[Back To The Top](#forum-go-project) 


## Installation

1. Clone the repository to your local machine.
2. Install Go if you haven't already.
3. Navigate to the project directory in your terminal.
4. Run `go run main.go` to start the web server.
5. Open your web browser and go to `http://localhost:8999` to access the application.<br><br>

[Back To The Top](#forum-go-project) 


## Requirements
Golang 1.23.0 or higher.<br><br>

[Back To The Top](#forum-go-project) 


## Learning Objectives
This project will help you learn about:

- Web development fundamentals:
  - HTML
  - HTTP
  - Sessions and cookies

- Database management:
  - SQL language basics
  - Database manipulation using SQLite

- Docker:
  - Setting up and using Docker
  - Containerizing applications
  - Creating Docker images

- User authentication:
  - Implementing user registration and login
  - Password encryption

- Web application features:
  - Creating and displaying posts and comments
  - Implementing like/dislike functionality
  - Filtering and categorizing content

- Error handling:
  - Managing website errors
  - Handling HTTP status codes
  - Addressing technical errors

- Best practices:
  - Following coding standards
  - Implementing unit testing (recommended)
<br><br>

By completing this project, you'll gain practical experience in full-stack web development, database management, and containerization, while also learning about user authentication and web application security.
<br><br>

[Back To The Top](#forum-go-project) 


## Last Updated
Last updated on August 15, 2026<br>

[Back To The Top](#forum-go-project) 


## Authors
### [Joon Kim](https://01.gritlab.ax/git/jkim)<br>
### [Mayuree Reunsati](https://01.gritlab.ax/git/mreunsat)
### [Geraldine Addamo](https://01.gritlab.ax/git/gaddamo)<br>
### [Sagynbek Osmonaliev](https://01.gritlab.ax/git/sosmonal)
### [Saddam Hussain](https://01.gritlab.ax/git/shussain)<br>

<br>

[Back To The Top](#forum-go-project) 

