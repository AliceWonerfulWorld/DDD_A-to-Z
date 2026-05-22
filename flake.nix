{
  description = "Lang War development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        pnpmPackage = pkgs.pnpm_10 or pkgs.pnpm;
        goPackage = pkgs.go_1_26 or pkgs.go;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.atlas
            pkgs.buf
            pkgs.docker
            pkgs.docker-compose
            pkgs.elixir
            pkgs.erlang
            pkgs.git
            pkgs.golangci-lint
            pkgs.nodejs_24
            pkgs.google-cloud-sdk
            pkgs.opentofu
            pnpmPackage
            goPackage
          ];
        };
      });
}
