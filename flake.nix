{
  description = "nux - A modern tmux session manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGoModule {
          pname = "nux";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;
          nativeBuildInputs = [ pkgs.installShellFiles ];
          ldflags = [
            "-s" "-w"
            "-X github.com/Drew-Daniels/nux/cmd.Version=${self.shortRev or "dev"}"
          ];
          postInstall = ''
            installShellCompletion --cmd nux \
              --bash <($out/bin/nux completions bash) \
              --fish <($out/bin/nux completions fish) \
              --zsh <($out/bin/nux completions zsh)
          '';
          meta.mainProgram = "nux";
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            golangci-lint
            gitleaks
            goreleaser
            lefthook
            hugo
            just
          ];
        };
      });
}
