require 'minitest/autorun'
require 'mocha/minitest'

# Assuming the DataModelService is in a relative path like 'internal/services/data_model_service.rb'
# Adjust the require path if needed based on your actual project structure
require_relative '../../internal/services/data_model_service'

class TestDataModelService < Minitest::Test

  def setup
    @db_backup_path = "/tmp/test_backups" # Use a temporary directory for tests
    Dir.mkdir(@db_backup_path) unless Dir.exist?(@db_backup_path)
    @data_model_service = DataModelService.new(@db_backup_path)
  end

  def teardown
    FileUtils.rm_rf(@db_backup_path) if Dir.exist?(@db_backup_path)
  end

  def test_backup_database_creates_backup_file
    # Mock the time to have a consistent backup name for the test.
    time_mock = mock()
    time_mock.expects(:strftime).with("%Y%m%d%H%M%S").returns("20240101000000")
    Time.expects(:now).returns(time_mock)

    # Call the backup method
    @data_model_service.backup_database

    # Assert that a file with the expected name exists in the backup directory
    expected_backup_file = File.join(@db_backup_path, "backup_20240101000000.dump")
    assert File.exist?(expected_backup_file), "Backup file should exist at #{expected_backup_file}"

    # Optionally, assert on the content of the file (if applicable)
    file_content = File.read(expected_backup_file)
    assert_equal "Simulated database backup data.", file_content, "Backup file should contain expected content"
  end

  def test_backup_database_handles_file_creation_error
    # Mock File.open to raise an error
    File.expects(:open).raises(StandardError, "Failed to create file")

    # Assert that the BackupDatabase method handles the error and returns an error
    assert_raises StandardError do
      @data_model_service.backup_database
    end
  end

end