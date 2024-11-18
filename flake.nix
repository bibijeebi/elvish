{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: {
    packages.x86_64-linux = {
      elvish = nixpkgs.legacyPackages.x86_64-linux.callPackage ./package.nix {};
    };
    devShells.x86_64-linux = nixpkgs.legacyPackages.x86_64-linux.mkShell {
      inputsFrom = [self.packages.x86_64-linux.elvish];
      buildInputs = [
        self.packages.x86_64-linux.elvish
        nixpkgs.legacyPackages.x86_64-linux.go
        nixpkgs.legacyPackages.x86_64-linux.gopls
        nixpkgs.legacyPackages.x86_64-linux.go-tools
        nixpkgs.legacyPackages.x86_64-linux.gofumpt
        nixpkgs.legacyPackages.x86_64-linux.gopls
        nixpkgs.legacyPackages.x86_64-linux.gopls-jsonrpc
        nixpkgs.legacyPackages.x86_64-linux.gopls-plugins
      ];
    };
  };
}
