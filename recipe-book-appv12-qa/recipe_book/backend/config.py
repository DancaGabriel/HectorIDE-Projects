import os

# Database configuration
DB_HOST = os.environ.get('DB_HOST', 'localhost')
DB_PORT = int(os.environ.get('DB_PORT', 3306))
DB_USER = os.environ.get('DB_USER', 'recipe_user')
DB_PASSWORD = os.environ.get('DB_PASSWORD', 'recipe_password')
DB_NAME = os.environ.get('DB_NAME', 'recipe_db')

# Flask application configuration
DEBUG = os.environ.get('DEBUG', True)  # Enable debugging mode
SECRET_KEY = os.environ.get('SECRET_KEY', 'super_secret_key') # Used for session management and CSRF protection

# Logging configuration (optional, can be further customized in logging_config.py)
LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')

# API settings
API_PREFIX = '/api'