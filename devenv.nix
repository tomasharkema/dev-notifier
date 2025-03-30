{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  # https://devenv.sh/basics/
  env.CGO_ENABLD = "1";

  # https://devenv.sh/packages/
  packages = with pkgs; [
    udev
    gcc
  ];

  # https://devenv.sh/languages/
  # languages.go.enable = true;
}
