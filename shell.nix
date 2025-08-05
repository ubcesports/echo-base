{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    # Go
    go

    # CLI tools
    sql-migrate
  ];
}
