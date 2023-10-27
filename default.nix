# build using 
# nix-build -E 'let pkgs = import <nixpkgs> { }; in pkgs.callPackage ./default.nix {}'
#
{ lib
, buildGoModule
, fetchFromGitHub
, nixosTests
, testers
, installShellFiles
}:
let
  version = "1.0.1";
  owner = "bezahl-online";
  pname = "rfidserver";
  repo = "RFID";
  rev = "v${version}";
  sha256 = "1ikyp1rxrl8lyfbll501f13yir1axighnr8x3ji3qzwin6i3w497";
in
buildGoModule {
  inherit version pname repo;

  src = fetchFromGitHub {
    owner = "bezahl-online";
    repo = repo;
    rev = "c4e06b0a215559e70ae77e0b3222428c2ea32e89";
    sha256 = "sha256-svWt11myC4QcHrO3uajau2FMC1iywA50aQzR2ma+dWk=";
  };
  # src = ../RFID/.;
 
  vendorSha256 = "sha256-GNZhzOp4orMBgJXIozOsswZXn7QR/ji1NFYwLFiw/3c=";

  buildPhase = ''
    runHook preBuild
    CGO_ENABLED=0 go build -o rfidserver .
    runHook postBuild
  '';

  installPhase = ''
    mkdir -p $out/bin
    mv rfidserver $out/bin
    cp localhost.crt localhost.key $out/bin
  '';

  meta = with lib; {
    homepage = "https://github.com/bezahl-online/RFID";
    description = "RFID server code";
    license = licenses.mit;
    maintainers = with maintainers; [ ralpheichelberger ];
  };
}

