let
  pkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/feda52be1d59.tar.gz") { };

in
pkgs.mkShell {
  buildInputs = with pkgs; [
    cargo
    rustc
    git-ignore
    cargo-watch
    rust-analyzer
    clippy
    darwin.apple_sdk.frameworks.Security
    gnumake
  ];

  RUST_SRC_PATH = "${pkgs.rust.packages.stable.rustPlatform.rustLibSrc}";
}
