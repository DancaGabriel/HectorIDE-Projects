from app.models.recipe import Recipe
from app.utils.db_utils import get_db_connection
from app.exceptions import RecipeNotFoundException

class RecipeService:
    """
    Service layer for handling recipe-related business logic.
    Interacts with the database and models to perform CRUD operations on recipes.
    """

    def __init__(self):
        """
        Initializes the RecipeService with a database connection.
        """
        self.conn = get_db_connection()
        self.cursor = self.conn.cursor()


    def get_all_recipes(self):
        """
        Retrieves all recipes from the database.

        Returns:
            list: A list of Recipe objects.
        """
        try:
            self.cursor.execute("SELECT id, title, ingredients, instructions FROM recipes")
            recipes = []
            for row in self.cursor.fetchall():
                recipe = Recipe(id=row[0], title=row[1], ingredients=row[2], instructions=row[3])
                recipes.append(recipe)
            return recipes
        except Exception as e:
            print(f"Error fetching all recipes: {e}")
            return [] # Or raise an exception, depending on error handling strategy


    def get_recipe_by_id(self, recipe_id):
        """
        Retrieves a specific recipe from the database by its ID.

        Args:
            recipe_id (int): The ID of the recipe to retrieve.

        Returns:
            Recipe: A Recipe object if found, None otherwise.

        Raises:
            RecipeNotFoundException: If no recipe is found with the given ID.
        """
        try:
            self.cursor.execute("SELECT id, title, ingredients, instructions FROM recipes WHERE id = %s", (recipe_id,))
            row = self.cursor.fetchone()

            if row:
                return Recipe(id=row[0], title=row[1], ingredients=row[2], instructions=row[3])
            else:
                raise RecipeNotFoundException(f"Recipe with id {recipe_id} not found")
        except RecipeNotFoundException as e:
            raise e  # Re-raise the custom exception
        except Exception as e:
            print(f"Error fetching recipe with id {recipe_id}: {e}")
            return None # Or raise an exception, depending on error handling strategy


    def create_recipe(self, title, ingredients, instructions):
        """
        Creates a new recipe in the database.

        Args:
            title (str): The title of the recipe.
            ingredients (str): The ingredients of the recipe (newline-separated).
            instructions (str): The instructions of the recipe.

        Returns:
            Recipe: The newly created Recipe object.
        """
        try:
            sql = "INSERT INTO recipes (title, ingredients, instructions) VALUES (%s, %s, %s)"
            self.cursor.execute(sql, (title, ingredients, instructions))
            self.conn.commit()
            new_recipe_id = self.cursor.lastrowid # get id of the last inserted row
            return self.get_recipe_by_id(new_recipe_id)
        except Exception as e:
            print(f"Error creating recipe: {e}")
            self.conn.rollback()
            return None  # Or raise an exception, depending on error handling strategy


    def update_recipe(self, recipe_id, title, ingredients, instructions):
        """
        Updates an existing recipe in the database.

        Args:
            recipe_id (int): The ID of the recipe to update.
            title (str): The new title of the recipe.
            ingredients (str): The new ingredients of the recipe (newline-separated).
            instructions (str): The new instructions of the recipe.

        Returns:
            Recipe: The updated Recipe object if successful, None otherwise.
        """
        try:
            sql = "UPDATE recipes SET title = %s, ingredients = %s, instructions = %s WHERE id = %s"
            self.cursor.execute(sql, (title, ingredients, instructions, recipe_id))
            self.conn.commit()
            return self.get_recipe_by_id(recipe_id)
        except RecipeNotFoundException as e:
             raise e # Re-raise the exception so controller can handle
        except Exception as e:
            print(f"Error updating recipe with id {recipe_id}: {e}")
            self.conn.rollback()
            return None  # Or raise an exception, depending on error handling strategy


    def delete_recipe(self, recipe_id):
        """
        Deletes a recipe from the database.

        Args:
            recipe_id (int): The ID of the recipe to delete.

        Returns:
            bool: True if the recipe was successfully deleted, False otherwise.
        """
        try:
            sql = "DELETE FROM recipes WHERE id = %s"
            self.cursor.execute(sql, (recipe_id,))
            self.conn.commit()
            return True
        except Exception as e:
            print(f"Error deleting recipe with id {recipe_id}: {e}")
            self.conn.rollback()
            return False  # Or raise an exception, depending on error handling strategy

    def __del__(self):
      """
      Ensures the database connection is closed when the service is destroyed.
      """
      if hasattr(self, 'conn') and self.conn:
        self.cursor.close()
        self.conn.close()