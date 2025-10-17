{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/7df7ff7d8e00218376575f0acdcc5d66741351ee.tar.gz") {config.allowUnfree = true;}}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    #app
    go
    golangci-lint
    go-task

    #ci
    act

    #infra
    openssl
    ansible
    vagrant
  ];

  shellHook = ''
    echo "Go version: $(go version)"
  '';
}