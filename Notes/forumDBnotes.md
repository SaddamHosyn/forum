## Database Schema

A database schema is a structured blueprint that defines how data is organized, stored, and related within a database. It serves as a logical representation of the entire database, outlining the tables, their attributes (columns), data types, constraints, relationships, and keys that maintain data integrity.
Key Components of a Database Schema
Tables: The fundamental units of a database where data is stored.
Attributes: The columns in each table that define the properties of the data.
Data Types: Specifications that dictate what kind of data can be stored in each attribute (e.g., integer, string).
Constraints: Rules applied to maintain data accuracy and integrity, including:
Primary Keys: Unique identifiers for records in a table.
Foreign Keys: Establish relationships between tables.
Unique Constraints: Ensure no duplicate entries exist in specified columns.
Not Null Constraints: Ensure that certain fields must contain values.
Types of Database Schemas
Conceptual Schema: Defines the entities, their attributes, and relationships without detailing how they will be physically stored.
Logical Schema: Describes the logical structure of the database, including tables, views, and relationships without focusing on physical storage details.
Physical Schema: Specifies how data is physically stored on storage devices, detailing aspects like data blocks and indexes.
Importance of Database Schemas
They guide users in performing operations such as querying, inserting, updating, and deleting data.
They help maintain consistency and integrity across the database by enforcing rules and relationships between different entities.
Database schemas are essential for efficient database design and management, enabling better performance and easier maintenance.



package main

import (
    "encoding/json"
    "net/http"
)

type Post struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

var posts []Post

func getPosts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func main() {
    http.HandleFunc("/posts", getPosts)
    http.ListenAndServe(":8080", nil)
}
