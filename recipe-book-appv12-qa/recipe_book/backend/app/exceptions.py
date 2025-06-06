class RecipeError(Exception):
    """Base class for recipe-related exceptions."""
    pass


class RecipeNotFound(RecipeError):
    """Exception raised when a recipe is not found."""
    def __init__(self, recipe_id):
        self.recipe_id = recipe_id
        self.message = f"Recipe with ID {recipe_id} not found."
        super().__init__(self.message)


class InvalidRecipeData(RecipeError):
    """Exception raised when recipe data is invalid."""
    def __init__(self, message="Invalid recipe data."):
        self.message = message
        super().__init__(self.message)


class DatabaseError(RecipeError):
    """Exception raised for database-related errors."""
    def __init__(self, message="A database error occurred."):
        self.message = message
        super().__init__(self.message)


class IngredientParsingError(RecipeError):
    """Exception raised when there's an issue parsing ingredients."""
    def __init__(self, message="Error parsing ingredients."):
        self.message = message
        super().__init__(self.message)