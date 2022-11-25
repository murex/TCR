# Fetch googletest

CPMFindPackage(
  NAME googletest
  GITHUB_REPOSITORY google/googletest
  GIT_TAG main
  OPTIONS "INSTALL_GTEST OFF" "gtest_force_shared_crt" "BUILD_GMOCK OFF"
)
