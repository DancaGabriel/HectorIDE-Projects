from flask import Flask
from app.routes.recipe_routes import recipe_bp
from app.utils.db_utils import db
from config import Config
import logging
from logging_config import setup_logging
from flask_cors import CORS

def create_app(config_class=Config):
    app = Flask(__name__)
    app.config.from_object(config_class)

    # Initialize Flask extensions here
    db.init_app(app)
    CORS(app)  # Enable CORS for all routes

    # Configure logging
    setup_logging(app)


    # Register blueprints
    app.register_blueprint(recipe_bp)

    with app.app_context():
        db.create_all()  # Create database tables

    return app


app = create_app()


@app.route('/')
def hello_world():
    app.logger.info('Hello world route accessed')
    return 'Hello, World!'

@app.errorhandler(500)
def internal_server_error(e):
    app.logger.error(f'Internal Server Error: {e}')
    return "Internal Server Error", 500


if __name__ == '__main__':
    app.run(debug=True)