#/bin/bash
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color
version=$(cat version)
echo "building version ${version}"
if [[ $(uname -m) != "armv7l" ]]; then echo -e "\n${RED}build need to be done on armv7 plattform${NC}\n"; exit 1; fi
(./build.sh && docker push www.greisslomat.at:5000/rfid:${version} && echo -e "\n${GREEN}push successfull${NC}\n") || echo -e "\n${RED}push failed${NC}\n"