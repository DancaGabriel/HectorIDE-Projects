require 'test/unit'
require 'dotenv'

# Load environment variables for testing
Dotenv.load

class TestConfig < Test::Unit::TestCase

  def test_environment_variables_are_loaded
    assert_not_nil(ENV['GIT_REPO_PATH'], "GIT_REPO_PATH should be defined")
    assert_not_nil(ENV['API_DOCUMENTATION_PATH'], "API_DOCUMENTATION_PATH should be defined")
    assert_not_nil(ENV['DATABASE_BACKUP_PATH'], "DATABASE_BACKUP_PATH should be defined")
  end

  def test_environment_variables_have_values
    assert(!ENV['GIT_REPO_PATH'].empty?, "GIT_REPO_PATH should not be empty")
    assert(!ENV['API_DOCUMENTATION_PATH'].empty?, "API_DOCUMENTATION_PATH should not be empty")
    assert(!ENV['DATABASE_BACKUP_PATH'].empty?, "DATABASE_BACKUP_PATH should not be empty")
  end

  # Add more tests as needed to validate the configuration loading process
  # such as type conversions or default value handling.
end