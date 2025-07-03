from sqlalchemy import Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class Recipe(Base):
    """
    Represents a recipe in the database.
    """
    __tablename__ = 'recipes'

    id = Column(Integer, primary_key=True, autoincrement=True)
    title = Column(String(255), nullable=False)
    ingredients = Column(String(4096), nullable=False)  # Store as newline-separated string
    instructions = Column(String(4096), nullable=False)

    def __repr__(self):
        return f"<Recipe(title='{self.title}')>"