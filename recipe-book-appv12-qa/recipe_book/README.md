# Recipe Book Application

A simple web application for managing recipes, built with Python (Flask) and MySQL.

## Project Overview

This application allows users to create, read, update, and delete (CRUD) their personal recipes. Each recipe consists of a title, a list of ingredients, and a set of instructions.

## Features

*   **Create Recipes:** Add new recipes to your collection.
*   **Read Recipes:** View details of your existing recipes.
*   **Update Recipes:** Modify existing recipes.
*   **Delete Recipes:** Remove recipes from your collection.
*   **RESTful API:**  Provides a REST API for managing recipes.

## Technologies Used

*   **Backend:**
    *   Python
    *   Flask (web framework)
    *   MySQL (database)
*   **Frontend:**
    *   HTML
    *   CSS
    *   JavaScript

## Project Structure

```
recipe_book/
├── backend/
│   ├── app/
│   │   ├── __init__.py
│   │   ├── models/
│   │   │   ├── __init__.py
│   │   │   └── recipe.py           # Recipe data model
│   │   ├── routes/
│   │   │   ├── __init__.py
│   │   │   └── recipe_routes.py    # API endpoints for recipes
│   │   ├── services/
│   │   │   ├── __init__.py
│   │   │   └── recipe_service.py   # Business logic for recipes
│   │   ├── utils/
│   │   │   ├── __init__.py
│   │   │   └── db_utils.py         # Database connection utilities
│   │   └── exceptions.py         # Custom exceptions
│   ├── config.py                 # Application configuration
│   ├── app.py                    # Main application instance
│   └── logging_config.py         # Logging setup
├── frontend/
│   ├── public/                    # Static assets directly served
│   │   ├── index.html            # Main HTML page
│   │   └── favicon.ico
│   ├── src/                       # Source files
│   │   ├── components/           # Reusable UI elements
│   │   │   └── RecipeList.js    # Recipe list component
│   │   ├── services/             # API interaction logic
│   │   │   └── api.js             # Function to call API
│   │   ├── App.js                 # Main frontend logic
│   │   └── styles.css             # Stylesheet
│   └── package.json               # Dependency declarations for frontend
├── requirements.txt              # Python dependencies
├── README.md
└── .env                        # Environment variables
```

## Setup and Installation

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd recipe_book
    ```

2.  **Backend Setup:**

    *   Navigate to the `backend` directory: `cd backend`
    *   Create a virtual environment (recommended):
        ```bash
        python3 -m venv venv
        source venv/bin/activate  # On Linux/macOS
        venv\Scripts\activate  # On Windows
        ```
    *   Install the required Python packages:
        ```bash
        pip install -r requirements.txt
        ```
    *   **Configure the database:**
        *   Create a MySQL database named `recipe_book`.
        *   Create a `.env` file in the `backend` directory with the following variables, replacing the values with your actual MySQL credentials:

            ```
            FLASK_APP=app.py
            FLASK_ENV=development # Or production
            DATABASE_HOST=localhost
            DATABASE_USER=your_mysql_user
            DATABASE_PASSWORD=your_mysql_password
            DATABASE_NAME=recipe_book
            DATABASE_PORT=3306  # Default MySQL port
            ```
    *   Run the application:
        ```bash
        flask run
        ```

3.  **Frontend Setup:**

    *   Navigate to the `frontend` directory: `cd ../frontend`
    *   Install the required JavaScript packages:

        ```bash
        npm install
        ```

    *   Start the frontend development server:

        ```bash
        npm start
        ```

    *   The frontend will typically be accessible at `http://localhost:3000`.

## API Endpoints

The backend provides the following REST API endpoints:

*   `GET /api/recipes`: Retrieve all recipes.
*   `GET /api/recipes/{id}`: Retrieve a specific recipe by ID.
*   `POST /api/recipes`: Create a new recipe.
*   `PUT /api/recipes/{id}`: Update an existing recipe.
*   `DELETE /api/recipes/{id}`: Delete a recipe.

## Database Setup (MySQL)

The application uses MySQL as its database. Make sure you have MySQL installed and running. The database schema will be created automatically by the application based on the `Recipe` data model defined in `backend/app/models/recipe.py`. The `.env` file must contain the correct database credentials.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues.

## License

[MIT](LICENSE)