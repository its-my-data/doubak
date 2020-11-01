# This command is used to compile protos into pb for being committed into source tree.
# Current build rules don't support compiling via bazel build directly:
# https://github.com/bazelbuild/rules_go/issues/512
protoc --proto_path=. --go_out=. --go_opt=paths=source_relative proto/*.proto

