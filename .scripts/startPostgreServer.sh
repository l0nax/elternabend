#!/bin/bash - 
#===============================================================================
#
#          FILE: startPostgreServer.sh
#
#         USAGE: ./startPostgreServer.sh
#
#   DESCRIPTION: This scripts starts an DEVELOP PostgreSQL Server in a Docker
#                container!
#                THIS IS ONLY FOR DEVELOPER!
#
#       OPTIONS: ---
#  REQUIREMENTS: docker
#          BUGS: ---
#         NOTES: ---
#        AUTHOR: Francesco Emanuel Bennici (l0nax), benniciemanuel78@gmail.com
#  ORGANIZATION: FABMation GmbH
#       CREATED: 01/15/2020 03:46:50 PM
#      REVISION:  001
#===============================================================================

set -o nounset                              # Treat unset variables as an error

DOCKER_CONTENT_TRUST=0

echo "****************************************"
echo "*             ATTENTION!               *"
echo "*   Please use this script only for    *"
echo "*   developing and not in production!  *"
echo "****************************************"
echo -e "\n\n"

set -x
sleep 1

docker pull postgres:10

docker run --rm -d \
        --name elternabend_db \
        -e POSTGRES_PASSWORD=postgres \
        -p 5432:5432 \
        postgres:10 \
        -c log_statement=all -c log_destination=stderr
