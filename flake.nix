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
        # TODO: 公式に提供されるまでは、Moonbit CLI の最新バイナリを直接ダウンロードして提供する。
        moonbitBinaries = {
          "aarch64-darwin" = {
            url = "https://cli.moonbitlang.com/binaries/latest/moonbit-darwin-aarch64.tar.gz";
            sha256 = "sha256-xv9K8uh5mV5eHooKh6DjnIIKv1ickhl1JV/YSUVW44Q=";
          };
          "aarch64-linux" = {
            url = "https://cli.moonbitlang.com/binaries/latest/moonbit-linux-aarch64.tar.gz";
            sha256 = "sha256-9nRdhueeaA2kMthHVdXKzVleSjrbu7ZP+uWI0Uzs5dY=";
          };
          "x86_64-linux" = {
            url = "https://cli.moonbitlang.com/binaries/latest/moonbit-linux-x86_64.tar.gz";
            sha256 = "sha256-ZwrAmmo6ZAC8lmeFbeZwXd7AqXzjvUwZxMOAtwfmhWw=";
          };
        };
        moonbit = pkgs.stdenv.mkDerivation {
          pname = "moonbit";
          version = "latest";
          src = pkgs.fetchurl {
            inherit (moonbitBinaries.${system}) url sha256;
          };
          sourceRoot = ".";
          nativeBuildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.autoPatchelfHook ];
          buildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.gcc.cc.lib ];
          installPhase = ''
            mkdir -p $out
            cp -r bin $out/bin
            cp -r lib $out/lib
            chmod -R +x $out/bin
          '';
          meta.platforms = builtins.attrNames moonbitBinaries;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.atlas
            pkgs.buf
            pkgs.dfmt
            pkgs.docker
            pkgs.docker-compose
            pkgs.dub
            pkgs.elixir
            pkgs.erlang
            pkgs.git
            pkgs.golangci-lint
            pkgs.ldc
            pkgs.nodejs_24
            pkgs.google-cloud-sdk
            pkgs.opentofu
            pkgs.php83
            pkgs.php83Packages.composer
            pkgs.github-linguist
            pnpmPackage
            goPackage
            moonbit
          ];
        };
      });
}
