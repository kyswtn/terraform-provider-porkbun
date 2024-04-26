{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs = { nixpkgs, ... }:
    let
      nixpkgsConfig = {
        allowUnfree = true;
      };
      mkInputs = system: {
        pkgs = import nixpkgs { inherit system; config = nixpkgsConfig; };
      };
      forAllSupportedSystems = fn:
        with nixpkgs.lib; genAttrs systems.flakeExposed (system: fn (mkInputs system));
    in
    {
      devShells = forAllSupportedSystems (inputs: with inputs; {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
            golangci-lint-langserver

            terraform
            terraform-ls
          ];
        };
      });
    };
}
