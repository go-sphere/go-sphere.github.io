package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var bindingLocationRE = regexp.MustCompile(`\bBINDING_LOCATION_[A-Z_]+\b`)

var knownBindingLocations = map[string]struct{}{
	"BINDING_LOCATION_UNSPECIFIED": {},
	"BINDING_LOCATION_QUERY":       {},
	"BINDING_LOCATION_URI":         {},
	"BINDING_LOCATION_JSON":        {},
	"BINDING_LOCATION_FORM":        {},
	"BINDING_LOCATION_HEADER":      {},
}

const githubAPIRepoLinkPrefix = "https://api.github.com/repos/" + "go-sphere/"

func TestContentDocsDoNotUseGitHubAPIRepoLinks(t *testing.T) {
	walkMarkdown(t, "content", func(path, text string) {
		if strings.Contains(text, githubAPIRepoLinkPrefix) {
			t.Errorf("%s contains a GitHub API repository link", path)
		}
	})
}

func TestContentDocsUseKnownBindingLocations(t *testing.T) {
	walkMarkdown(t, "content", func(path, text string) {
		for _, token := range bindingLocationRE.FindAllString(text, -1) {
			if _, ok := knownBindingLocations[token]; !ok {
				t.Errorf("%s uses unknown binding enum %s", path, token)
			}
		}
	})
}

func TestContentDocsProtoExamplesCompile(t *testing.T) {
	protoc, err := exec.LookPath("protoc")
	if err != nil {
		t.Skip("protoc is not installed")
	}

	fixtureDir := writeProtoFixtures(t)
	snippetDir := t.TempDir()
	var compiled int

	walkMarkdown(t, "content", func(path, text string) {
		for _, block := range protobufBlocks(text) {
			if !isCompleteProtoExample(block.source) {
				continue
			}

			compiled++
			protoFile := filepath.Join(snippetDir, fmt.Sprintf("snippet_%03d.proto", compiled))
			if err := os.WriteFile(protoFile, []byte(block.source), 0o644); err != nil {
				t.Fatalf("write %s: %v", protoFile, err)
			}

			outFile := filepath.Join(t.TempDir(), "snippet.pb")
			cmd := exec.Command(
				protoc,
				"--proto_path="+snippetDir,
				"--proto_path="+fixtureDir,
				"--descriptor_set_out="+outFile,
				protoFile,
			)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("%s:%d proto example does not compile:\n%s", path, block.line, strings.TrimSpace(string(output)))
			}
		}
	})

	if compiled == 0 {
		t.Fatal("no complete protobuf examples were found")
	}
}

type protobufBlock struct {
	line   int
	source string
}

func walkMarkdown(t *testing.T, root string, fn func(path, text string)) {
	t.Helper()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fn(path, string(content))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func protobufBlocks(text string) []protobufBlock {
	var blocks []protobufBlock
	var block strings.Builder
	var startLine int
	inProto := false

	for i, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			fence := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			if inProto {
				blocks = append(blocks, protobufBlock{line: startLine, source: block.String()})
				block.Reset()
				inProto = false
				continue
			}

			inProto = fence == "proto" || fence == "protobuf"
			startLine = i + 1
			continue
		}

		if inProto {
			block.WriteString(line)
			block.WriteByte('\n')
		}
	}

	return blocks
}

func isCompleteProtoExample(source string) bool {
	return strings.Contains(source, `syntax = "proto3";`) && strings.Contains(source, "package ")
}

func writeProtoFixtures(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	fixtures := map[string]string{
		"buf/validate/validate.proto": `syntax = "proto3";

package buf.validate;

import "google/protobuf/descriptor.proto";

message FieldConstraints {
  bool required = 1;
  Int64Rules int64 = 2;
  Int32Rules int32 = 3;
}

message Int64Rules {
  int64 gt = 1;
}

message Int32Rules {
  int32 gt = 1;
}

extend google.protobuf.FieldOptions {
  FieldConstraints field = 1159;
}
`,
		"google/api/annotations.proto": `syntax = "proto3";

package google.api;

import "google/protobuf/descriptor.proto";

message HttpRule {
  string get = 2;
  string put = 3;
  string post = 4;
  string delete = 5;
  string patch = 6;
  string body = 7;
  repeated HttpRule additional_bindings = 11;
  string response_body = 12;
}

extend google.protobuf.MethodOptions {
  HttpRule http = 72295728;
}
`,
		"shared/v1/user.proto": `syntax = "proto3";

package shared.v1;

message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
}
`,
		"sphere/binding/binding.proto": `syntax = "proto3";

package sphere.binding;

import "google/protobuf/descriptor.proto";

enum BindingLocation {
  BINDING_LOCATION_UNSPECIFIED = 0;
  BINDING_LOCATION_QUERY = 1;
  BINDING_LOCATION_URI = 2;
  BINDING_LOCATION_JSON = 3;
  BINDING_LOCATION_FORM = 4;
  BINDING_LOCATION_HEADER = 5;
}

extend google.protobuf.MessageOptions {
  optional BindingLocation default_location = 136655300;
  repeated string default_auto_tags = 136655301;
}

extend google.protobuf.OneofOptions {
  optional BindingLocation default_oneof_location = 136655310;
  repeated string default_oneof_auto_tags = 136655311;
}

extend google.protobuf.FieldOptions {
  optional BindingLocation location = 136655320;
  repeated string tags = 136655321;
  repeated string auto_tags = 136655322;
}
`,
		"sphere/errors/errors.proto": `syntax = "proto3";

package sphere.errors;

import "google/protobuf/descriptor.proto";

message Error {
  int32 status = 1;
  string reason = 2;
  string message = 3;
}

extend google.protobuf.EnumOptions {
  int32 default_status = 18534200;
}

extend google.protobuf.EnumValueOptions {
  Error options = 18534210;
}
`,
		"sphere/options/options.proto": `syntax = "proto3";

package sphere.options;

import "google/protobuf/descriptor.proto";

message KeyValuePair {
  string key = 1;
  oneof value {
    bool flag = 2;
    string text = 3;
    int64 number = 4;
  }
  map<string, string> extra = 5;
}

extend google.protobuf.MethodOptions {
  repeated KeyValuePair options = 501319300;
}
`,
	}

	for path, content := range fixtures {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("create fixture dir for %s: %v", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("write fixture %s: %v", path, err)
		}
	}

	return root
}
