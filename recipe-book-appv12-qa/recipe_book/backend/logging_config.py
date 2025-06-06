import logging
import logging.config
import os

def setup_logging(default_path='logging.yaml', default_level=logging.INFO, env_key='LOG_CFG'):
    """Setup logging configuration"""
    path = default_path
    value = os.getenv(env_key, None)
    if value:
        path = value
    if os.path.exists(path):
        import yaml
        with open(path, 'rt') as f:
            try:
                config = yaml.safe_load(f.read())
                logging.config.dictConfig(config)
            except Exception as e:
                print(f"Error loading logging configuration: {e}")
                logging.basicConfig(level=default_level)
    else:
        logging.basicConfig(level=default_level)
        print('Failed to load configuration file. Using default configs')

if __name__ == '__main__':
    # Example usage:  Create a dummy logging.yaml file if it doesn't exist
    # to demonstrate default behavior.
    if not os.path.exists('logging.yaml'):
        with open('logging.yaml', 'w') as f:
            f.write("""
version: 1
formatters:
  simple:
    format: '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
handlers:
  console:
    class: logging.StreamHandler
    level: DEBUG
    formatter: simple
    stream: ext://sys.stdout
root:
  level: INFO
  handlers: [console]
""")
    setup_logging()
    logger = logging.getLogger(__name__)
    logger.info("This is an informational message.")
    logger.debug("This is a debug message.") #won't show up unless logging level is changed above.