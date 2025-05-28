require 'minitest/autorun'
require 'mocha/minitest'
require_relative '../../internal/services/git_service'

class GitServiceTest < Minitest::Test
  def setup
    @repo_path = "test_repo"
    @git_service = GitService.new(@repo_path)
  end

  def teardown
    # Clean up the test repository after each test.  This should be safe because
    # we initialize a repo in a temporary directory each test.
    FileUtils.rm_rf(@repo_path) if Dir.exist?(@repo_path)
  end

  def test_initialize_git_repo_creates_repo_if_not_exists
    refute Dir.exist?(@repo_path + "/.git")
    @git_service.InitializeGitRepo
    assert Dir.exist?(@repo_path + "/.git")
  end

  def test_initialize_git_repo_does_not_error_if_repo_exists
    Dir.mkdir(@repo_path)
    Dir.mkdir(@repo_path + "/.git")
    assert_nil @git_service.InitializeGitRepo #Should not raise an exception
  end

  def test_commit_and_push_adds_commits_and_pushes
    Dir.mkdir(@repo_path)
    @git_service.InitializeGitRepo

    # Create a dummy file for the commit
    File.write(File.join(@repo_path, "test_file.txt"), "test content")

    # Mock the execution of git commands
    git_add_mock = mock
    git_add_mock.expects(:run).returns(true)
    git_commit_mock = mock
    git_commit_mock.expects(:run).returns(true)
    git_push_mock = mock
    git_push_mock.expects(:run).returns(true) #Return true to avoid printing warning in the service.

    Open3.expects(:capture2e).with("git", "-C", @repo_path, "add", ".").returns([nil, git_add_mock])
    Open3.expects(:capture2e).with("git", "-C", @repo_path, "commit", "-m", "Test Commit").returns([nil, git_commit_mock])
    Open3.expects(:capture2e).with("git", "-C", @repo_path, "push").returns([nil, git_push_mock])

    @git_service.CommitAndPush("Test Commit")
  end

  def test_commit_and_push_handles_push_failure
    Dir.mkdir(@repo_path)
    @git_service.InitializeGitRepo

    # Create a dummy file for the commit
    File.write(File.join(@repo_path, "test_file.txt"), "test content")

    # Mock the execution of git commands
    git_add_mock = mock
    git_add_mock.expects(:run).returns(true)
    git_commit_mock = mock
    git_commit_mock.expects(:run).returns(true)
    git_push_mock = mock
    git_push_mock.expects(:run).raises("Push Failed")

    Open3.expects(:capture2e).with("git", "-C", @repo_path, "add", ".").returns([nil, git_add_mock])
    Open3.expects(:capture2e).with("git", "-C", @repo_path, "commit", "-m", "Test Commit").returns([nil, git_commit_mock])
    Open3.expects(:capture2e).with("git", "-C", @repo_path, "push").returns([nil, git_push_mock])

    @git_service.CommitAndPush("Test Commit")
  end

  def test_commit_and_push_handles_add_failure
      Dir.mkdir(@repo_path)
      @git_service.InitializeGitRepo

      # Mock the execution of git commands
      git_add_mock = mock
      git_add_mock.expects(:run).raises("Add Failed")

      Open3.expects(:capture2e).with("git", "-C", @repo_path, "add", ".").returns([nil, git_add_mock])

      assert_raises(RuntimeError) { @git_service.CommitAndPush("Test Commit") }
  end

  def test_commit_and_push_handles_commit_failure
    Dir.mkdir(@repo_path)
    @git_service.InitializeGitRepo

    # Create a dummy file for the commit
    File.write(File.join(@repo_path, "test_file.txt"), "test content")

    # Mock the execution of git commands
    git_add_mock = mock
    git_add_mock.expects(:run).returns(true)
    git_commit_mock = mock
    git_commit_mock.expects(:run).raises("Commit Failed")

    Open3.expects(:capture2e).with("git", "-C", @repo_path, "add", ".").returns([nil, git_add_mock])
    Open3.expects(:capture2e).with("git", "-C", @repo_path, "commit", "-m", "Test Commit").returns([nil, git_commit_mock])

    assert_raises(RuntimeError) { @git_service.CommitAndPush("Test Commit") }
  end
end