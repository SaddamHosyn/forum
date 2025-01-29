# Forum Project Guideline
<br>
This guideline will help you complete your forum project in one week. Each day is broken into specific tasks with estimated time allocations.
<br>

## Day 1: Project Setup and Database Design

### Tasks:
- **Set up the project environment** (2 hours)
  - Install necessary tools (Go, Docker, SQLite).
  - Create a new Go project.

- **Design the database schema** (4 hours)
  - Define tables: `Users`, `Posts`, `Comments`, `Categories`, `Likes/Dislikes`.
  - Create an entity relationship diagram to visualize relationships.
  - Example table structure:<br>

    - **Users**: `id` (UUID), `email`, `username`, `password` (hashed).
    - **Posts**: `id`, `user_id` (FK), `title`, `content`, `category_id` (FK), `created_at`.
    - **Comments**: `id`, `post_id` (FK), `user_id` (FK), `content`, `created_at`.
    - **Categories**: `id`, `name`.
    - **Likes/Dislikes**: `id`, `user_id` (FK), `post_id` (FK), `type` (like/dislike).
<br>
<br>


## Day 2: Implement Database Operations
### Tasks:
- **Create SQLite database and tables** (3 hours)
  - Write SQL queries to create tables based on your design.

- **Implement basic CRUD operations** (5 hours)
  - Create functions for inserting users, posts, comments.
  - Implement SELECT queries to retrieve data.
  <br>
<br>


## Day 3: User Authentication
### Tasks:
- **Develop user registration functionality** (4 hours)
  - Create registration endpoint that checks for existing emails and hashes passwords.<br>

- **Implement login session management** (4 hours)
  - Create login endpoint using cookies for session management.
  - Set cookie expiration and handle session validation.<br>


## Day 4: Post and Comment Functionality
### Tasks:
- **Enable users to create posts and comments** (4 hours)
  - Develop endpoints for creating posts and comments.
  
- **Implement visibility rules for registered vs. non-registered users** (4 hours)
  - Ensure only registered users can create posts/comments but all can view them.


## Day 5: Likes and Dislikes Feature
### Tasks:
- **Add liking/disliking functionality** (4 hours)
  - Create endpoints to like/dislike posts and comments.
  
- **Display like/dislike counts publicly** (4 hours)
  - Update the post and comment retrieval functions to include like/dislike counts.



## Day 6: Filtering Mechanism
### Tasks:
- **Implement filtering options for posts** (5 hours)
  - Allow filtering by categories and user-specific liked posts.
  
- **Test filtering functionality thoroughly** (3 hours)
  - Ensure filters work correctly under various scenarios.



## Day 7: Testing and Deployment
### Tasks:
- **Write unit tests for all functionalities** (4 hours)
  - Ensure that all endpoints work as expected and handle errors gracefully.

- 
