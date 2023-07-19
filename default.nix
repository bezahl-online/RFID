{ lib
, buildGoModule
, nixosTests
, testers
, installShellFiles
}:
let
  version = "1.0.1";
  owner = "bezahl-online";
  repo = "RFID";
  rev = "v${version}";
  sha256 = "1ikyp1rxrl8lyfbll501f13yir1axighnr8x3ji3qzwin6i3w497";
in
buildGoModule {
  pname = "rfidserver";
  inherit version;

  src = ./.;
 
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
    maintainers = with maintainers; [ /* list of maintainers here */ ];
  };
}

