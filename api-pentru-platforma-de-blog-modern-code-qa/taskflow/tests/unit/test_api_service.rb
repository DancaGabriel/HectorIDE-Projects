require 'minitest/autorun'
require 'mocha/minitest'
require_relative '../../internal/services/api_service'

class TestAPIService < Minitest::Test
  def setup
    @doc_path = 'test_api_doc.txt'
    @api_service = APIService.new(@doc_path)
  end

  def teardown
    File.delete(@doc_path) if File.exist?(@doc_path)
  end

  def test_update_api_documentation_success
    new_content = "Updated API Documentation"
    assert_nil @api_service.UpdateAPIDocumentation(new_content) # Should return nil on success
    assert_equal new_content, File.read(@doc_path)
  end

  def test_update_api_documentation_failure
    # Mock os to raise an error when writing the file
    File.expects(:write).with(@doc_path, "some content").raises(Errno::EACCES.new("Permission denied"))
    api_service = APIService.new(@doc_path)

    error = assert_raises(RuntimeError) do
      api_service.UpdateAPIDocumentation("some content")
    end

    assert_match "failed to update API documentation", error.message
    assert_match "Permission denied", error.message
  end

  def test_new_api_service
      api_service = APIService.new("some path")
      assert_equal "some path", api_service.doc_path
  end
end