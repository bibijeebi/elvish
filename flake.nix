{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    packages.${system} = {
      elvish = pkgs.callPackage ./package.nix {};
      default = self.packages.${system}.elvish;
    };

    devShells.${system}.default = pkgs.mkShell {
      inputsFrom = [self.packages.${system}.elvish];
      buildInputs = with pkgs; [
        go
        gopls
        go-tools
        gofumpt
        make
      ];

      # Add these environment variables for Go development
      shellHook = ''
        export GOPATH="$PWD/.go"
        export GOBIN="$GOPATH/bin"
        export PATH="$GOBIN:$PATH"

        # For plugin development (optional)
        export CGO_ENABLED=1
      '';
    };
  };
}
