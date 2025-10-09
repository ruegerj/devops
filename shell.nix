{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/7df7ff7d8e00218376575f0acdcc5d66741351ee.tar.gz") {}}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    golangci-lint
    go-task
  ];

  shellHook = ''
    echo "Go version: $(go version)"
  '';
}