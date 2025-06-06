from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
import logging
import os

# Initialize extensions
db = SQLAlchemy()
migrate = Migrate()

def create_app(config_filename=None):
    """
    Application factory function.  Initializes the Flask application and
    registers blueprints.
    """
    app = Flask(__name__)

    # Load default configuration
    app.config.from_object('config.Config')

    # Override with environment variables if set.  This is safer than
    # including secret keys directly in config.py
    for key, value in os.environ.items():
        if key.startswith('RECIPE_BOOK_'):  # Custom prefix to avoid collisions
            app.config[key[12:]] = value  # Remove prefix

    # Optionally load a configuration file.  This allows setting the
    # configuration programmatically in tests, for example.
    if config_filename:
        app.config.from_pyfile(config_filename)

    # Initialize extensions
    db.init_app(app)
    migrate.init_app(app, db)

    # Configure logging
    logging.basicConfig(level=logging.INFO)
    app.logger.setLevel(logging.INFO)

    # Register blueprints
    from .routes import recipe_routes
    app.register_blueprint(recipe_routes.recipe_bp)

    app.logger.info("Application initialized")

    return app

# Import models (needed for Alembic to detect them)
from . import models